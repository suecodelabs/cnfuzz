package file

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"

	"github.com/suecodelabs/cnfuzz/src/log"
	"github.com/suecodelabs/cnfuzz/src/model"
)

type containerImageJsonRepository struct {
	filePath string
}

func CreateContainerImageRepository() *containerImageJsonRepository {
	jsonFile := "./container_images.json"
	if _, err := os.Stat(jsonFile); errors.Is(err, os.ErrNotExist) {
		// File doesn't exist yet
		file, err := os.Create(jsonFile)
		if err != nil {
			log.L().Fatalf("failed to create file (%s) that is used by the ContainerImageRepository to store fuzzed image information", jsonFile)
		}
		_, err = file.WriteString("[]")
		if err != nil {
			log.L().Panic("failed to initialize Json file for persistence")
			return nil
		}
		_ = file.Close()
	}

	return &containerImageJsonRepository{
		filePath: jsonFile,
	}
}

func (repo containerImageJsonRepository) GetContainerImages() ([]model.ContainerImage, error) {
	fileBytes, err := os.ReadFile(repo.filePath)
	if err != nil {
		return nil, err
	}

	var images []model.ContainerImage
	err = json.Unmarshal(fileBytes, &images)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (repo containerImageJsonRepository) FindContainerImageByHash(hash string) (containerImage model.ContainerImage, found bool, err error) {
	images, err := repo.GetContainerImages()
	if err != nil {
		return model.ContainerImage{}, false, err
	}
	for _, image := range images {
		if image.Hash == hash {
			return image, true, nil
		}
	}

	return model.ContainerImage{}, false, nil
}

func (repo containerImageJsonRepository) CreateContainerImage(image model.ContainerImage) error {
	valErr := image.Verify()
	if valErr != nil {
		return valErr
	}

	images, err := repo.GetContainerImages()
	if err != nil {
		return err
	}
	images = append(images, image)

	// Read existing file
	jsonBytes, err := json.Marshal(images)
	if err != nil {
		return err
	}

	// Open file and truncate content
	file, err := os.OpenFile(repo.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	// Write new content to file
	writer := bufio.NewWriter(file)
	_, err = writer.Write(jsonBytes)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	// Close file
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (repo containerImageJsonRepository) UpdateContainerImage(image model.ContainerImage) error {
	return nil
}
