package utils

import (
	"fmt"

	"github.com/ldicarlo/legifrss/server/models"
)

func ExtractAndConvertDILA(input []models.JOContainerResult) (result []models.LegifranceElement) {

	for _, jorfContent := range input {
		for _, item := range jorfContent.Items {
			for _, step := range item.Container.Structure.Contents {
				result = tranformHierarchyStep(step, []string{}, result)
			}
		}
	}
	return
}

func transformSummary(summary models.Summary, categories []string) models.LegifranceElement {
	return models.LegifranceElement{
		Id:          summary.Id,
		Title:       summary.Id,
		Link:        "https://www.legifrance.gouv.fr/jorf/id/" + summary.Id,
		Description: summary.Title,
		Category:    categories,
		Author:      summary.Emitter,
		Nature:      summary.Nature,
	}
}

func tranformHierarchyStep(
	hs models.HierarchyStep,
	categories []string,
	result []models.LegifranceElement,
) []models.LegifranceElement {

	for _, item := range hs.Summaries {
		result = append(result, transformSummary(item, append(categories, hs.Title)))
	}
	for _, nextHs := range hs.Step {
		result = tranformHierarchyStep(nextHs, append(categories, hs.Title), result)
	}
	return result
}

func ErrCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic("Panic.")
	}
}

func ErrCheckStr(str string) {
	if str != "" {
		fmt.Println(str)
		panic("Panic.")
	}
}
