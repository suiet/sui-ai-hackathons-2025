package llm

//
//func TestValidateNFTImageURI(t *testing.T) {
//	isSafe, invalidateCategories, err := ValidateNFTImage("title", "description",
//		"https://static.wixstatic.com/media/33e670_8d4bd0510f4f43dd855353779441aec8~mv2.jpg")
//	if err != nil {
//		t.Error(err)
//	}
//
//	if isSafe {
//		t.Error("Image should be not safe")
//	}
//
//	if len(invalidateCategories) == 0 {
//		t.Log("No invalidate categories", invalidateCategories)
//	}
//
//	isSafe, invalidateCategories, err = ValidateNFTImage("title", "description",
//		"https://bq7bkvdje7gvgmv66hrxdy7wx5h5ggtrrnmt66rdkkehb64rvz3q.arweave.net/DD4VVGknzVMyvvHjceP2v0_TGnGLWT96I1KIcPuRrnc")
//	if err != nil {
//		t.Error(err)
//	}
//
//	if !isSafe {
//		t.Error("Image should be safe")
//	}
//
//	if len(invalidateCategories) == 0 {
//		t.Log("No invalidate categories", invalidateCategories)
//	}
//}
