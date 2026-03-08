package service

import (
	"os"
	"path/filepath"

	authDto "teacher-os-api/internal/modules/auth/dto"
	"teacher-os-api/internal/shared/errs"

	"github.com/lukasjarosch/go-docx"
)

func (s *ExportService) ExportLessonPlanDOCX(
	currentUser *authDto.UserResponse,
	planID string,
) ([]byte, string, error) {
	doc, err := s.GetLessonPlanPreview(currentUser, planID)
	if err != nil {
		return nil, "", err
	}

	templatePath := filepath.Join("templates", "lesson-plan-th.docx")

	wordDoc, err := docx.Open(templatePath)
	if err != nil {
		return nil, "", errs.Internal(err)
	}

	data := buildTemplateData(doc)

	replaceMap := docx.PlaceholderMap{}
	for key, value := range data {
		replaceMap[key] = value
	}

	if err := wordDoc.ReplaceAll(replaceMap); err != nil {
		return nil, "", errs.Internal(err)
	}

	if err := os.MkdirAll("tmp", 0o755); err != nil {
		return nil, "", errs.Internal(err)
	}

	outputPath := filepath.Join("tmp", "lesson-plan-"+doc.PlanID+".docx")

	if err := wordDoc.WriteToFile(outputPath); err != nil {
		return nil, "", errs.Internal(err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, "", errs.Internal(err)
	}

	filename := "lesson-plan-" + doc.PlanID + ".docx"
	return content, filename, nil
}
