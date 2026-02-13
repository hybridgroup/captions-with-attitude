package main

import (
	"errors"
	"fmt"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
)

// VLM is a Vision Language Model (VLM).
type VLM struct {
	TextModelFilename      string
	ProjectorModelFilename string

	TextModel        llama.Model
	Sampler          llama.Sampler
	ModelContext     llama.Context
	ProjectorContext mtmd.Context

	template string
}

// NewVLM creates a new VLM.
func NewVLM(model, projector string) *VLM {
	return &VLM{
		TextModelFilename:      model,
		ProjectorModelFilename: projector,
	}
}

// Close closes the VLM.
func (m *VLM) Close() {
	if m.ProjectorContext != 0 {
		mtmd.Free(m.ProjectorContext)

	}

	if m.ModelContext != 0 {
		llama.Free(m.ModelContext)
	}
}

// Init initializes the VLM.
func (m *VLM) Init() error {
	var err error
	m.TextModel, err = llama.ModelLoadFromFile(m.TextModelFilename, llama.ModelDefaultParams())
	if err != nil {
		return fmt.Errorf("unable to load text model: %w", err)
	}

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 2048

	m.ModelContext, err = llama.InitFromModel(m.TextModel, ctxParams)
	if err != nil {
		return fmt.Errorf("unable to initialize model context: %w", err)
	}

	m.template = llama.ModelChatTemplate(m.TextModel, "")

	m.Sampler = llama.NewSampler(m.TextModel, llama.DefaultSamplers, llama.DefaultSamplerParams())

	mtmdCtxParams := mtmd.ContextParamsDefault()
	m.ProjectorContext, err = mtmd.InitFromFile(m.ProjectorModelFilename, m.TextModel, mtmdCtxParams)
	if err != nil {
		return fmt.Errorf("unable to initialize projector context: %w", err)
	}

	return nil
}

// ChatTemplate applies the model's chat template to the given messages.
func (m *VLM) ChatTemplate(messages []llama.ChatMessage, add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(m.template, messages, add, buf)
	result := string(buf[:len])

	return result
}

// Tokenize tokenizes the input text and image bitmap into output chunks.
func (m *VLM) Tokenize(input *mtmd.InputText, bitmap mtmd.Bitmap, output mtmd.InputChunks) (err error) {
	if res := mtmd.Tokenize(m.ProjectorContext, output, input, []mtmd.Bitmap{bitmap}); res != 0 {
		err = fmt.Errorf("unable to tokenize: %d", res)
	}
	return
}

// Results generates text results from the given input chunks.
func (m *VLM) Results(output mtmd.InputChunks) (string, error) {
	var n llama.Pos
	nBatch := llama.NBatch(m.ModelContext)

	if res := mtmd.HelperEvalChunks(m.ProjectorContext, m.ModelContext, output, 1, 0, int32(nBatch), true, &n); res != 0 {
		return "", errors.New("unable to evaluate chunks")
	}

	vocab := llama.ModelGetVocab(m.TextModel)
	results := ""

	for i := 0; i < int(nBatch); i++ {
		token := llama.SamplerSample(m.Sampler, m.ModelContext, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 128)
		len := llama.TokenToPiece(vocab, token, buf, 0, true)
		results += string(buf[:len])

		batch := llama.BatchGetOne([]llama.Token{token})
		batch.Pos = &n

		llama.Decode(m.ModelContext, batch)
		n++
	}

	m.Clear()

	return results, nil
}

// Clear clears the context memory.
func (m *VLM) Clear() {
	llama.Synchronize(m.ModelContext)
	mem, err := llama.GetMemory(m.ModelContext)
	if err != nil {
		fmt.Println("unable to get memory:", err)
		return
	}
	llama.MemoryClear(mem, true)
}
