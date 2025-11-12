package handlers

import()

type dirReq struct{
	Path string `json:"path"`
	Name string `json:"name"`
}

type renameReq struct{
	Path    string `json:"path"`
	Name    string `json:"name"`
	NewName string `json:"newName"`
}

type deleteReq struct{
	Path string `json:"path"`
	Name string `json:"name"`
}