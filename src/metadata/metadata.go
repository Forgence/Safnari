package metadata

func ExtractMetadata(path string, mimeType string) map[string]interface{} {
    metadata := make(map[string]interface{})

    switch mimeType {
    case "image/jpeg", "image/png":
        meta := extractImageMetadata(path)
        for k, v := range meta {
            metadata[k] = v
        }
    case "application/pdf":
        meta := extractPDFMetadata(path)
        for k, v := range meta {
            metadata[k] = v
        }
    case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
        meta := extractDOCXMetadata(path)
        for k, v := range meta {
            metadata[k] = v
        }
    default:
        // Unsupported MIME type for metadata extraction
    }

    return metadata
}

func extractImageMetadata(path string) map[string]interface{} {
    // Implement EXIF extraction using go-exif library
    return nil
}

func extractPDFMetadata(path string) map[string]interface{} {
    // Implement PDF metadata extraction using unipdf library
    return nil
}

func extractDOCXMetadata(path string) map[string]interface{} {
    // Implement DOCX metadata extraction using gooxml library
    return nil
}