package services

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

const docAnalysisErr = "document analysis failed"

type IdentityFields struct {
	FirstName  string
	MiddleName string
	LastName   string
}

type TextractService interface {
	ExtractIDContent(ctx context.Context, key string) (map[string]string, error)
}

type textractService struct {
	bucket         string
	textractClient *textract.Client
	logger         *slog.Logger
}

type blockMap map[string]types.Block
type wordsMap map[string]string

func (t *textractService) generateBlockMaps(docBlocks []types.Block) (blockMap, blockMap, wordsMap) {
	blocks := make(blockMap)
	keyBlocks := make(blockMap)
	words := make(wordsMap)

	for _, block := range docBlocks {
		blocks[aws.ToString(block.Id)] = block

		if block.BlockType == types.BlockTypeKeyValueSet && block.EntityTypes[0] == types.EntityTypeKey {
			keyBlocks[aws.ToString(block.Id)] = block
		}

		if block.BlockType == types.BlockTypeWord {
			words[aws.ToString(block.Id)] = aws.ToString(block.Text)
		}
	}

	return blocks, keyBlocks, words
}

func (t *textractService) findKey(rels []types.Relationship, words wordsMap) string {
	var key string

	for _, rel := range rels {
		if rel.Type != types.RelationshipTypeChild {
			continue
		}

		for _, id := range rel.Ids {
			key += words[id]
		}
	}

	return strings.ToLower(key)
}

func (t *textractService) findValue(rels []types.Relationship, blocks blockMap, words wordsMap) string {
	var word []string

	for _, rel := range rels {
		if rel.Type != types.RelationshipTypeValue {
			continue
		}

		for _, id := range rel.Ids {
			for _, childId := range blocks[id].Relationships[0].Ids {
				word = append(word, words[childId])
			}
		}
	}

	return strings.Join(word, " ")
}

func (t *textractService) findKeysWords(blocks blockMap, keyBlocks blockMap, words wordsMap) wordsMap {
	keyWords := make(wordsMap, len(keyBlocks))

	for _, keyBlock := range keyBlocks {
		key := t.findKey(keyBlock.Relationships, words)
		value := t.findValue(keyBlock.Relationships, blocks, words)
		keyWords[key] = value
	}

	return keyWords
}

func (t *textractService) ExtractIDContent(ctx context.Context, key string) (map[string]string, error) {
	output, err := t.textractClient.AnalyzeDocument(ctx, &textract.AnalyzeDocumentInput{
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeForms,
		},
		Document: &types.Document{
			S3Object: &types.S3Object{
				Bucket: aws.String(t.bucket),
				Name:   aws.String(key),
			},
		},
	})

	if err != nil {
		t.logger.ErrorContext(ctx, docAnalysisErr, "error", err)
		return nil, err
	}

	blocks, keyBlocks, words := t.generateBlockMaps(output.Blocks)
	fields := t.findKeysWords(blocks, keyBlocks, words)

	return fields, nil
}

func NewTextractService(bucket string, cfg aws.Config, logger *slog.Logger) TextractService {
	textactClient := textract.NewFromConfig(cfg)

	return &textractService{
		bucket:         bucket,
		textractClient: textactClient,
		logger:         logger,
	}
}
