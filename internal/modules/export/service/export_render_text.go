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
		"subjectId":             fallbackString(doc.SubjectID),
		"learningUnitId":        fallbackString(doc.LearningUnitID),
		"lessonNo":              intToString(doc.LessonNo),
		"lessonTitle":           fallbackString(doc.LessonTitle),
		"lessonHours":           intToString(doc.LessonHours),
		"schoolName":            fallbackString(doc.SchoolName),
		"teacherName":           fallbackString(doc.TeacherName),
		"subjectGroup":          fallbackString(doc.SubjectGroup),
		"gradeLevel":            fallbackString(doc.GradeLevel),
		"semester":              fallbackString(doc.Semester),
		"academicYear":          fallbackString(doc.AcademicYear),
		"planId":                fallbackString(doc.PlanID),
		"generatedAt":           time.Now().Format("02/01/2006 15:04"),
		"basicInfoText":         renderBasicInfoText(doc.BasicInfo),
		"objectiveText":         renderObjectiveText(doc.Objective),
		"standardIndicatorText": renderStandardIndicatorText(doc.StandardIndicator),
		"activityText":          renderActivityText(doc.Activity),
		"mediaMaterialText":     renderMediaMaterialText(doc.MediaMaterial),
		"assessmentText":        renderAssessmentText(doc.Assessment),
		"homeworkText":          renderHomeworkText(doc.Homework),
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

func renderBasicInfoText(data map[string]interface{}) string {
	var lines []string

	lines = appendIfValue(lines, "ชื่อแผนการจัดการเรียนรู้", data["ชื่อแผนการจัดการเรียนรู้"])
	lines = appendIfValue(lines, "รายวิชา", data["รายวิชา"])
	lines = appendIfValue(lines, "รหัสวิชา", data["รหัสวิชา"])
	lines = appendIfValue(lines, "ระดับชั้น", data["ระดับชั้น"])
	lines = appendIfValue(lines, "ภาคเรียน", data["ภาคเรียน"])
	lines = appendIfValue(lines, "ปีการศึกษา", data["ปีการศึกษา"])
	lines = appendIfValue(lines, "หน่วยการเรียนรู้", data["หน่วยการเรียนรู้"])
	lines = appendIfValue(lines, "เรื่อง", data["เรื่อง"])
	lines = appendIfValue(lines, "เวลาเรียน", data["เวลาเรียน"])

	return fallbackSection(lines)
}

func renderObjectiveText(data map[string]interface{}) string {
	var lines []string

	lines = appendTextBlock(lines, "สาระสำคัญ", data["สาระสำคัญ"])
	lines = appendListBlock(lines, "จุดประสงค์การเรียนรู้", data["จุดประสงค์การเรียนรู้"])

	return fallbackSection(lines)
}

func renderStandardIndicatorText(data map[string]interface{}) string {
	var lines []string

	lines = appendListBlock(lines, "มาตรฐานการเรียนรู้", data["มาตรฐานการเรียนรู้"])
	lines = appendListBlock(lines, "ตัวชี้วัด", data["ตัวชี้วัด"])
	lines = appendListBlock(lines, "สมรรถนะสำคัญของผู้เรียน", data["สมรรถนะสำคัญของผู้เรียน"])
	lines = appendListBlock(lines, "คุณลักษณะอันพึงประสงค์", data["คุณลักษณะอันพึงประสงค์"])

	return fallbackSection(lines)
}

func renderActivityText(data map[string]interface{}) string {
	var lines []string

	lines = appendTextBlock(lines, "ขั้นนำเข้าสู่บทเรียน", data["ขั้นนำเข้าสู่บทเรียน"])
	lines = appendTextBlock(lines, "ขั้นจัดกิจกรรมการเรียนรู้", data["ขั้นจัดกิจกรรมการเรียนรู้"])
	lines = appendTextBlock(lines, "ขั้นสรุป", data["ขั้นสรุป"])
	lines = appendListBlock(lines, "กิจกรรมย่อย", data["กิจกรรมย่อย"])

	return fallbackSection(lines)
}

func renderMediaMaterialText(data map[string]interface{}) string {
	var lines []string

	lines = appendListBlock(lines, "สื่อการเรียนรู้", data["สื่อการเรียนรู้"])
	lines = appendListBlock(lines, "อุปกรณ์", data["อุปกรณ์"])
	lines = appendListBlock(lines, "แหล่งการเรียนรู้", data["แหล่งการเรียนรู้"])

	return fallbackSection(lines)
}

func renderAssessmentText(data map[string]interface{}) string {
	var lines []string

	lines = appendListBlock(lines, "วิธีการวัดและประเมินผล", data["วิธีการวัดและประเมินผล"])
	lines = appendListBlock(lines, "เครื่องมือ", data["เครื่องมือ"])
	lines = appendListBlock(lines, "เกณฑ์การประเมิน", data["เกณฑ์การประเมิน"])

	return fallbackSection(lines)
}

func renderHomeworkText(data map[string]interface{}) string {
	var lines []string

	lines = appendTextBlock(lines, "ชิ้นงานหรือภาระงาน", data["ชิ้นงานหรือภาระงาน"])
	lines = appendTextBlock(lines, "งานที่มอบหมาย", data["งานที่มอบหมาย"])
	lines = appendTextBlock(lines, "บันทึกหลังสอน", data["บันทึกหลังสอน"])
	lines = appendTextBlock(lines, "หมายเหตุ", data["หมายเหตุ"])

	return fallbackSection(lines)
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

func intToString(value int) string {
	if value <= 0 {
		return "-"
	}
	return fmt.Sprintf("%d", value)
}

func appendIfValue(lines []string, label string, value interface{}) []string {
	text := normalizeText(value)
	if text == "" {
		return lines
	}
	return append(lines, fmt.Sprintf("%s: %s", label, text))
}

func appendTextBlock(lines []string, label string, value interface{}) []string {
	text := normalizeText(value)
	if text == "" {
		return lines
	}

	if len(lines) > 0 {
		lines = append(lines, "")
	}

	lines = append(lines, label+":")
	lines = append(lines, text)
	return lines
}

func appendListBlock(lines []string, label string, value interface{}) []string {
	items := normalizeList(value)
	if len(items) == 0 {
		return lines
	}

	if len(lines) > 0 {
		lines = append(lines, "")
	}

	lines = append(lines, label+":")
	for _, item := range items {
		lines = append(lines, "- "+item)
	}

	return lines
}

func normalizeText(value interface{}) string {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return ""
	}
}

func normalizeList(value interface{}) []string {
	raw, ok := value.([]interface{})
	if !ok {
		if typed, ok := value.([]string); ok {
			var result []string
			for _, item := range typed {
				item = strings.TrimSpace(item)
				if item != "" {
					result = append(result, item)
				}
			}
			return result
		}
		return nil
	}

	var result []string
	for _, item := range raw {
		text := strings.TrimSpace(fmt.Sprint(item))
		if text != "" {
			result = append(result, text)
		}
	}

	return result
}

func fallbackSection(lines []string) string {
	if len(lines) == 0 {
		return "-"
	}
	return strings.Join(lines, "\n")
}
