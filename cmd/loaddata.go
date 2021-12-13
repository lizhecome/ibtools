package cmd

import (
	"ibtools_server/models"
	"ibtools_server/util"
	"log"
	"path/filepath"
)

//LoadData 加载信号用卡数据
func LoadData(filename, filetype, configBackend string) error {
	fullfilename, _ := filepath.Abs(filename)
	if !util.Exist(fullfilename) {
		log.Printf("数据文件不存在！%s", fullfilename)
	}

	_, db, _, _, _, err := initConfigDB(true, false, configBackend)
	if err != nil {
		return err
	}

	prj := new(models.Project)
	//prj.Code = uuid.New().String()
	prj.IsTemplate = 1
	prj.Type = "IPO"
	prj.PicUrl = "https://www.baidu.com/s?wd=%E7%99%BE%E5%BA%A6%E7%83%AD%E6%90%9C&sa=ire_dl_gh_logo_texing&rsv_dl=igh_logo_pcs"
	prj.DDModelList = make([]models.DDModel, 0)
	prj.Tilte = "IPO模板"

	model := new(models.DDModel)
	model.Order = 1
	model.Title = "公司基本情况"

	item := new(models.DDItem)
	item.Order = 1
	item.ReviewMethod = "核查方法（大段文字）"
	item.Title = "改革与设立情况"
	item.Status = "未分配"
	item.FilePointer = "1-2-1"

	// item.BaseInfoSchema = "{}"
	// item.BaseInfos = "{}"

	question := new(models.ReviewQuestion)
	question.Order = 1
	question.QuestionDescription = "问题描述"
	question.QuestionName = "问题名称"
	question.Solution = "解决方案"
	question.LawSupport = "法律支持大段文字"

	additionitem := new(models.DDItem)
	additionitem.Order = 1
	additionitem.ReviewMethod = "核查方法（大段文字）"
	additionitem.Title = "附加核查项目"
	additionitem.Status = "未分配"
	additionitem.FilePointer = "1-2-1"

	// additionitem.BaseInfoSchema = "{}"
	// additionitem.BaseInfos = "{}"

	question.AdditionDDItem = make([]models.DDItem, 0)
	question.AdditionDDItem = append(question.AdditionDDItem, *additionitem)

	model.Questions = make([]models.ReviewQuestion, 0)
	model.Questions = append(model.Questions, *question)

	model.DDItems = make([]models.DDItem, 0)
	model.DDItems = append(model.DDItems, *item)

	prj.DDModelList = make([]models.DDModel, 0)
	prj.DDModelList = append(prj.DDModelList, *model)
	db.Create(&prj)

	// jsonPessoal, _ := json.Marshal(prj)
	// fmt.Fprintf(os.Stdout, "%s", jsonPessoal)

	// if fileObj, err := os.Open(fullfilename); err == nil {
	// 	defer fileObj.Close()
	// 	//一个文件对象本身是实现了io.Reader的 使用bufio.NewReader去初始化一个Reader对象，存在buffer中的，读取一次就会被清空
	// 	reader := bufio.NewReader(fileObj)
	// 	if filetype == "card" {
	// 		for {
	// 			line, _, err := reader.ReadLine()

	// 			if err == io.EOF {
	// 				break
	// 			}
	// 			var cc models.CreditCard
	// 			if err := json.Unmarshal(line, &cc); err != nil {
	// 				fmt.Println(err)
	// 			} else {
	// 				fmt.Println(cc)
	// 				imagefilename := filepath.Join(filepath.Dir(fullfilename), "ccimage", fmt.Sprintf("%s.png", cc.ID))
	// 				_, cdnpath, _ := util.UpLoadFile(imagefilename, "ccimage")
	// 				cc.ImageLink = cdnpath

	// 				existcc := new(models.CreditCard)
	// 				notFound := db.Where("id = ?", cc.ID).First(existcc).RecordNotFound()
	// 				if notFound {
	// 					if err := db.Create(&cc).Error; err != nil {
	// 						fmt.Println(err)
	// 					}
	// 				} else {
	// 					cc.ID = existcc.ID
	// 					if err := db.Save(&cc).Error; err != nil {
	// 						fmt.Println(err)
	// 					}
	// 				}

	// 			}
	// 			//fmt.Printf("%s \n", line)
	// 		}
	// 	}
	// 	if filetype == "googlerelation" {
	// 		for {
	// 			line, _, err := reader.ReadLine()

	// 			if err == io.EOF {
	// 				break
	// 			}
	// 			var rr models.GooglePlaceType2PurchaseType
	// 			if err := json.Unmarshal(line, &rr); err != nil {
	// 				fmt.Println(err)
	// 			} else {
	// 				fmt.Println(rr)
	// 				if err := db.Create(&rr).Error; err != nil {
	// 					fmt.Println(err)
	// 				}
	// 			}
	// 			//fmt.Printf("%s \n", line)
	// 		}
	// 	}
	// 	if filetype == "yelprelation" {
	// 		for {
	// 			line, _, err := reader.ReadLine()

	// 			if err == io.EOF {
	// 				break
	// 			}
	// 			var rr models.YelpPlaceType2PurchaseType
	// 			if err := json.Unmarshal(line, &rr); err != nil {
	// 				fmt.Println(err)
	// 			} else {
	// 				fmt.Println(rr)
	// 				if err := db.Create(&rr).Error; err != nil {
	// 					fmt.Println(err)
	// 				}
	// 			}
	// 			//fmt.Printf("%s \n", line)
	// 		}
	// 	}
	// 	if filetype == "recommendation" {
	// 		for {
	// 			line, _, err := reader.ReadLine()

	// 			if err == io.EOF {
	// 				break
	// 			}
	// 			var rr models.Recommendation
	// 			if err := json.Unmarshal(line, &rr); err != nil {
	// 				fmt.Println(err)
	// 			} else {
	// 				fmt.Println(rr)
	// 				if err := db.Create(&rr).Error; err != nil {
	// 					fmt.Println(err)
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if filetype == "yelp" {
	// 		data, err := ioutil.ReadFile(fullfilename)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		var yelps []models.YelpCategory
	// 		err = json.Unmarshal(data, &yelps)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		// var level1 []string
	// 		// for _, v := range yelps {
	// 		// 	if len(v.Parents) == 0 {
	// 		// 		level1 = append(level1, v.Alias)
	// 		// 	}
	// 		// }

	// 		// for _, v := range yelps {
	// 		// 	if len(v.Parents) != 0 {
	// 		// 		var islevel2 bool
	// 		// 		islevel2 = false
	// 		// 		for _, vv := range v.Parents {
	// 		// 			if ok, _ := util.Contain(vv, level1); ok {
	// 		// 				islevel2 = true
	// 		// 				break
	// 		// 			}
	// 		// 		}
	// 		// 		if islevel2 {
	// 		// 			fmt.Printf("####%s$$$$\n", v.Alias)
	// 		// 		}
	// 		// 	}
	// 		// }

	// 		for _, v := range yelps {
	// 			if err := db.Create(&v).Error; err != nil {
	// 				fmt.Println(err)
	// 			}
	// 		}
	// 	}
	// }
	return nil

}
