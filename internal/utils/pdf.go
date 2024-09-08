package utils

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

type PDFBuilder struct {
	pdf *gofpdf.Fpdf
}

// NewPDFBuilder initializes a new PDF builder
func NewPDFBuilder() *PDFBuilder {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	return &PDFBuilder{pdf: pdf}
}

// AddTitle adds a title to the PDF
func (b *PDFBuilder) AddTitle(title string) *PDFBuilder {
	b.pdf.SetFont("Arial", "B", 16)
	b.pdf.Cell(40, 10, title)
	b.pdf.Ln(12)
	b.pdf.SetFont("Arial", "", 12)
	return b
}

// AddField adds a field label and value to the PDF
func (b *PDFBuilder) AddField(label, value string) *PDFBuilder {
	b.pdf.Cell(40, 10, fmt.Sprintf("%s: %s", label, value))
	b.pdf.Ln(10)
	return b
}

// AddSectionTitle adds a section title
func (b *PDFBuilder) AddSectionTitle(title string) *PDFBuilder {
	b.pdf.SetFont("Arial", "B", 14)
	b.pdf.Cell(40, 10, title)
	b.pdf.Ln(10)
	b.pdf.SetFont("Arial", "", 12)
	return b
}

// Output generates the PDF and returns it as a bytes.Buffer
func (b *PDFBuilder) Output() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := b.pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
