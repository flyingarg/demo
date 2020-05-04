package data

import (
	"demo/internal/env"
	"demo/internal/logger"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Item struct {
	ID string `json:"id" xml:"id"`
	Name string `json:"name" xml:"name"`
	Manufacturer string `json:"manufacturer" xml:"manufacturer"`
	Brand string `json:"brand" xml:"brand"`
	Category string `json:"category" xml:"category"`
	Images []Image `json:"images" xml:"images"`
}

type Image struct {
	URL string `json:"url" xml:"url"`
	ID string `json:"id" xml:"id"`
	Status bool `json:"status" xml:"status"`
	Error string `json:"error" xml:"error"`
	Location string `json:"location" xml:"location"`
}

type Request struct {
	ID string `json:"id" xml:"id"`
	Items []Item `json:"items" xml:"items"`
	Type string
}


//ProcessRequest takes ths post data as []bytes.
//This data is checked for being a json or an xml and then unmarshalled to data.Request struct.
//The Request then checks the total number of items that a request can contain.
//We then check, the plan's request quota has been exceeded. If yes, the user is intimated as such.(this is preliminary check)
//After the sanity checks, the application, would download the images on each item. If an item's image download fails,
//the image's status field is marked as false. If the images status field is set to false, the item's status would also
//be false.
func ProcessRequest(data []byte, username string, plan string) (Request, error) {
	logger.Log.Sugar().Infof("processing request %s", string(data))
	var request Request
	var err error
	dataType := ""
	if json.Valid(data){
		err = json.Unmarshal(data, &request)
		dataType = "json"
	}else{
		err = xml.Unmarshal(data, &request)
		dataType = "xml"
	}
	request.Type = dataType
	if err != nil {
		return Request{}, errors.New("neither a valid xml or json")
	}

	if len(request.Items) > 50 {
		return Request{}, errors.New("exceeded permissible number of items")
	}

	if ok, err2 := WithinRateLimits(username, plan); !ok || err2 != nil {
		return Request{}, err2
	}

	//Adding ids so that the status can be tracked in the future
	allImages := make(map[string]Image)
	requestId, _ := uuid.NewRandom()
	request.ID = requestId.String()
	i := 0

	var newItems []Item
	for _, item := range request.Items{
		item.ID = request.ID + "_" + strconv.Itoa(i)
		i += 1
		j := 0
		var newImages []Image
		for _, image := range item.Images {
			image.ID = item.ID + "_" + strconv.Itoa(j)
			allImages[image.ID] = image
			newImages = append(newImages, Image{
				URL:      image.URL,
				ID:       image.ID,
				Status:   image.Status,
			})
			j += 1
		}
		newItems = append(newItems, Item{
			ID:           item.ID,
			Name:         item.Name,
			Manufacturer: item.Manufacturer,
			Brand:        item.Brand,
			Category:     item.Category,
			Images:       newImages,
		})
	}

	newReq := Request{
		ID:    requestId.String(),
		Items:  newItems,
	}

	var wg sync.WaitGroup
	c := make(chan Image)
	fetchAllImages(c, &wg, allImages)
	go func() {
		(&wg).Wait()
		close(c)
	}()

	for image := range c {
		allImages[image.ID] = image
	}


	var newRItems []Item
	for _, item := range newReq.Items {
		var newRImages []Image
		for _, image := range item.Images {
			image = allImages[image.ID]
			newRImages = append(newRImages, Image{
				URL:      image.URL,
				ID:       image.ID,
				Status:   image.Status,
				Location: image.Location,
				Error: image.Error,
			})
		}
		newRItems = append(newRItems, Item{
			ID:           item.ID,
			Name:         item.Name,
			Manufacturer: item.Manufacturer,
			Brand:        item.Brand,
			Category:     item.Category,
			Images:       newRImages,
		})
	}
	newRReq := Request{
		ID:    requestId.String(),
		Items:  newRItems,
		Type:  dataType,
	}

	if ok, err2 := HonorRateLimitPostProessing(username, plan, request, dataType, string(data));  !ok || err2 != nil{
		return Request{}, err2
	}
	return newRReq, nil
}

//fetchAllImages Routine to invoke multiple getImage go routines
func fetchAllImages(out chan Image, wg *sync.WaitGroup, allImages map[string]Image){
	logger.Log.Sugar().Infof("total images to download %d", len(allImages))
	for _, v := range allImages {
		wg.Add(1)
		go getImage(out, wg, v)
	}
}

//getImage fetched the data from the provided url.
//If it succeeds to fetch the file, save it to a file system, it returns on chanel of type image where the image
//struct's Status flag is set to true. Also, Location of the image is populated.
//
//If it fails to save the image, it set image.Status as false and image.Error is populated.
func getImage(out chan Image, wg *sync.WaitGroup, image Image){
	logger.Log.Sugar().Infof("downloading image %s", image.URL)
	defer wg.Done()
	localLocation := ""
	resp, err := http.Get(image.URL)
	if err != nil {
		logger.Log.Sugar().Errorf("failed to download image with error %e", err)
		image.Status = false
		image.Error = err.Error()
		out <- image
		return
	}else{
		defer resp.Body.Close()
		localLocation = env.GetEnv("FS_IMAGE_LOC").(string) + "/test_" + image.ID
		logger.Log.Sugar().Infof("writing image to location %s", localLocation)
		o, err3 := os.Create(localLocation)
		if err3 != nil {
			logger.Log.Sugar().Errorf("failed to write image with error %e", err)
			image.Status = false
			image.Error = err3.Error()
			out <- image
			return
		}
		defer o.Close()
		_, err2 := io.Copy(o, resp.Body)
		if err2 != nil {
			image.Status = true
			image.Error = err2.Error()
			image.Location = localLocation
			out <- image
			return
		}
	}

	image.Error = ""
	image.Status = true
	image.Location = localLocation
	out <- image
}
