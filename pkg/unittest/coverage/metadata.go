package coverage

import (
	"fmt"
	"regexp"
	"strings"
)

type CoverageReportMetadata struct {
	Templates map[string][]CoverageReportMetadataTemplate
	TotalCharts int
	TotalManifests int
	TotalDocuments int
}

type CoverageReportMetadataTemplate struct {
	Path string
	FullPath string
	Documents []string
}

func NewCoverageReportMetadata() *CoverageReportMetadata {
	templates := make(map[string][]CoverageReportMetadataTemplate)
	return &CoverageReportMetadata{
		Templates: templates,
	}
}

var splitterPattern = regexp.MustCompile("(?m:^---$)")

func(c *CoverageReportMetadata) AppendTemplate(chartPath string, fileName string, content string) {
	parts := splitterPattern.Split(string(content), -1)
	documents := make([]string, 0)
	for _, part := range parts {
		if len(strings.TrimSpace(part)) > 0 {
			documents = append(documents, strings.TrimSpace(part))
			c.TotalDocuments++
		}
	}
	c.Templates[chartPath] = append(c.Templates[chartPath], CoverageReportMetadataTemplate{
		Path: fileName,
		FullPath: fmt.Sprintf("%s/%s", chartPath, fileName),
		Documents: documents,
	})
	c.TotalManifests++
	c.TotalCharts = len(c.Templates)
}

func(c *CoverageReportMetadata) OutElements() {
	// for _, elements := range c.Templates {
	// 	fmt.Println(elements)
	// }
	
	fmt.Println("Total Manifests", c.TotalManifests)
    fmt.Println("Total Documents", c.TotalDocuments)
}
