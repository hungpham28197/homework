package controller

import (
	"os"
)

const FOLDER_UPLOAD = "./uploads/"

func RemoveAvatarFile(avatarFileName string) {
	os.Remove(FOLDER_UPLOAD + avatarFileName)
}
