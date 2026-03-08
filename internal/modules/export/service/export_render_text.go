package service

import (
	"fmt"
	"strings"
	"time"

	exportDto "teacher-os-api/internal/modules/export/dto"
)

func buildTemplateData(doc *exportDto.LessonPlanDocument) map[string]string {
	return map[string]string{
		"documentTitle":         "แผนการจัดการเรียนรู้",
		"title":                 fallbackString(doc.Title),
		"schoolName":            fallbackString(doc.SchoolName),
		"teacherName":           fallbackString(doc.TeacherName),
		"subjectGroup":          fallbackString(doc.SubjectGroup),
		"gradeLevel":            fallbackString(doc.GradeLevel),
		"semester":              fallbackString(doc.Semester),
		"academicYear":          fallbackString(doc.AcademicYear),
		"planId":                fallbackString(doc.PlanID),
		"generatedAt":           time.Now().Format("02/01/2006 15:04"),
		"basicInfoText":         renderSectionText(doc.BasicInfo),
		"objectiveText":         renderSectionText(doc.Objective),
		"standardIndicatorText": renderSectionText(doc.StandardIndicator),
		"activityText":          renderSectionText(doc.Activity),
		"mediaMaterialText":     renderSectionText(doc.MediaMaterial),
		"assessmentText":        renderSectionText(doc.Assessment),
		"homeworkText":          renderSectionText(doc.Homework),
	}
}

func renderSectionText(data map[string]interface{}) string {
	if len(data) == 0 {
		return "-"
	}

	var lines []string
	for key, value := range data {
		lines = append(lines, renderField(key, value)...)
	}

	return strings.Join(lines, "\n")
}

func renderField(label string, value interface{}) []string {
	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return []string{fmt.Sprintf("%s: -", label)}
		}
		return []string{fmt.Sprintf("%s: %s", label, v)}

	case []interface{}:
		if len(v) == 0 {
			return []string{fmt.Sprintf("%s: -", label)}
		}

		lines := []string{fmt.Sprintf("%s:", label)}
		for _, item := range v {
			lines = append(lines, fmt.Sprintf("- %v", item))
		}
		return lines

	default:
		return []string{fmt.Sprintf("%s: %v", label, v)}
	}
}

func fallbackString(value string) string {
	if strings.TrimSpace(value) == "" {
		return "-"
	}
	return value
}
