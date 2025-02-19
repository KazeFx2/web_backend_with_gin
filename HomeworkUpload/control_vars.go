package HomeworkUpload

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"main/Config"
	"main/Fs"
	"main/Logger"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type ExpControl struct {
	Exp     int       `yaml:"exp"`
	Start   time.Time `yaml:"start"`
	End     time.Time `yaml:"end"`
	ReStart time.Time `yaml:"re_start"`
	ReEnd   time.Time `yaml:"re_end"`
}

type HomeworkVar struct {
	Control struct {
		Auth        string       `yaml:"authorization"`
		MaxExp      int          `yaml:"max_exp"`
		SaveDir     string       `yaml:"save_dir"`
		GroupExps   []int        `yaml:"grp_exps"`
		ExpControls []ExpControl `yaml:"exp_controls"`
	} `yaml:"control"`
	Students []Student `yaml:"students"`
}

var HwControl = HomeworkVar{
	Control: struct {
		Auth        string       `yaml:"authorization"`
		MaxExp      int          `yaml:"max_exp"`
		SaveDir     string       `yaml:"save_dir"`
		GroupExps   []int        `yaml:"grp_exps"`
		ExpControls []ExpControl `yaml:"exp_controls"`
	}{
		Auth:        "",
		MaxExp:      8,
		SaveDir:     "./fileSave",
		ExpControls: []ExpControl{},
	},
	Students: make([]Student, 0),
}

type Student struct {
	StudentName string `yaml:"student_name"`
	StudentNum  string `yaml:"student_num"`
	StudentGrp  int    `yaml:"student_grp"`
}

type FileRecord struct {
	Stu        Student
	LastUpload time.Time
	Hash       string
	Uploaded   bool
}

type DynamicExp struct {
	Exp     int
	Start   time.Time
	End     time.Time
	ReStart time.Time
	ReEnd   time.Time
	Exist   bool
	Records []FileRecord
}

var Starts = make([]time.Time, 0)
var Ends = make([]time.Time, 0)
var ReStarts = make([]time.Time, 0)
var ReEnds = make([]time.Time, 0)
var ExpAva = make([]bool, 0)

var Students = make([]Student, 0)

var StudentNumMap = make(map[string]int)

var StudentNameMap = make(map[string]int)

var GroupStu = make(map[int][]int)

var DynamicExpInfo = make([]DynamicExp, 0)

func InitLoad(filePath string) error {
	st, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = Config.WriteYaml(filePath, &HwControl)
			if err != nil {
				Logger.LogE("init yaml '%s' failed: %v", filePath, err)
				return err
			}
			Logger.LogI("init yaml '%s' success", filePath)
			return nil
		} else {
			Logger.LogE("open file '%s' failed: %v", filePath, err)
			return err
		}
	}
	if st.IsDir() {
		Logger.LogE("'%s' is a directory!", filePath)
		return errors.New("config file can not be a directory")
	}
	err = Config.LoadYaml(filePath, &HwControl)
	if err != nil {
		Logger.LogE("load config from '%s' failed: %v", filePath, err)
		return err
	}
	return nil
}

func LoadStudents(filePath string) error {
	_, name := filepath.Split(filePath)
	postFix := filepath.Ext(name)
	if postFix == ".xls" {
		err := Config.TransformXls2Xlsx(filePath, strings.ReplaceAll(filePath, postFix, ".xlsx"))
		if err != nil {
			return err
		}
		filePath = strings.ReplaceAll(filePath, postFix, ".xlsx")
	}
	rows, err := Config.LoadExcel(filePath)
	if err != nil {
		Logger.LogE("load file '%s' failed: %v", filePath, err)
		return err
	}
	index := 0
	for _, col := range rows {
		var st Student
		flag := false
		for _, cell := range col {
			if cell == "" {
				continue
			}
			if !flag {
				st.StudentNum = strings.ReplaceAll(cell, " ", "")
				flag = true
			} else {
				st.StudentName = strings.ReplaceAll(cell, " ", "")
				break
			}
		}
		Students = append(Students, st)
		HwControl.Students = append(HwControl.Students, st)
		StudentNumMap[st.StudentNum] = index
		StudentNameMap[st.StudentName] = index
		index++
	}
	return nil
}

func InitStartsEnds() {
	for i := 1; i <= HwControl.Control.MaxExp; i++ {
		flag := false
		for _, Ec := range HwControl.Control.ExpControls {
			if Ec.Exp == i {
				Starts = append(Starts, Ec.Start)
				Ends = append(Ends, Ec.End)
				ReStarts = append(ReStarts, Ec.ReStart)
				ReEnds = append(ReEnds, Ec.ReEnd)
				ExpAva = append(ExpAva, true)
				flag = true
				break
			}
		}
		if !flag {
			Starts = append(Starts, time.Now())
			Ends = append(Ends, time.Now())
			ReStarts = append(ReStarts, time.Now())
			ReEnds = append(ReEnds, time.Now())
			ExpAva = append(ExpAva, false)
		}
	}
}

func InitGroups() {
	for _, i := range HwControl.Students {
		arr, exist := GroupStu[i.StudentGrp]
		if !exist {
			GroupStu[i.StudentGrp] = []int{StudentNumMap[i.StudentNum]}
		} else {
			GroupStu[i.StudentGrp] = append(arr, StudentNumMap[i.StudentNum])
		}
	}
}

func InitRecords(rc *[]FileRecord) {
	for _, st := range Students {
		*rc = append(*rc, FileRecord{
			Stu:        st,
			LastUpload: time.Now(),
			Hash:       "",
			Uploaded:   false,
		})
	}
}

func InitGrpRcs(rc *[]FileRecord) {
	for i := 1; i <= len(GroupStu); i++ {
		*rc = append(*rc, FileRecord{
			Stu:        Student{StudentGrp: i + 1, StudentNum: "", StudentName: ""},
			LastUpload: time.Now(),
			Hash:       "",
			Uploaded:   false,
		})
	}
}

func IsGroupExp(expNum int) bool {
	for _, i := range HwControl.Control.GroupExps {
		if i == expNum {
			return true
		}
	}
	return false
}

func InitDynamicInfo() error {
	InitStartsEnds()
	InitGroups()
	for i := 1; i <= HwControl.Control.MaxExp; i++ {
		var De = DynamicExp{
			Exp:     i,
			Start:   Starts[i-1],
			End:     Ends[i-1],
			ReStart: ReStarts[i-1],
			ReEnd:   ReEnds[i-1],
			Exist:   ExpAva[i-1],
			Records: make([]FileRecord, 0),
		}
		if IsGroupExp(i) {
			InitGrpRcs(&De.Records)
		} else {
			InitRecords(&De.Records)
		}
		targetPath := path.Join(HwControl.Control.SaveDir, fmt.Sprintf("%d", i))
		if !Fs.DirAva(targetPath) {
			err := os.Mkdir(targetPath, 0644)
			if err != nil {
				Logger.LogE("failed in 'mkdir' of dir '%s': %v", targetPath, err)
				return err
			}
		}
		Et, err := os.ReadDir(targetPath)
		if err != nil {
			Logger.LogE("failed in read dir '%s': %v", targetPath, err)
			return err
		}
		for _, e := range Et {
			if e.IsDir() {
				continue
			}
			fullName := path.Join(targetPath, e.Name())
			fileName := e.Name()
			groupExp := IsGroupExp(i)
			name, StudentNum, StudentName, StudentGrp := "", "", "", -1
			if !groupExp {
				name = strings.ReplaceAll(strings.ReplaceAll(fileName, filepath.Ext(fullName), ""), fmt.Sprintf("第%d次实验", i), "")
				StudentNum = name[:9]
				StudentName = name[9:]
			} else {
				name = strings.ReplaceAll(fileName, filepath.Ext(fullName), "")
				_, err := fmt.Sscanf(name, "第%d组大作业", &StudentGrp)
				if err != nil {
					return err
				}
			}
			info, err := e.Info()
			if err != nil {
				Logger.LogE("failed to read info of '%s': %v", fullName, err)
				return err
			}
			Time := info.ModTime()
			ior, err := os.ReadFile(fullName)
			if err != nil {
				Logger.LogE("failed to open file '%s': %v", fullName, err)
				return err
			}
			hash := sha256.New()
			if _, err := io.Copy(hash, bytes.NewReader(ior)); err != nil {
				Logger.LogE("failed in calc hash of file '%s': %v", fullName, err)
				return err
			}
			hashBytes := hash.Sum(nil)
			Hash := hex.EncodeToString(hashBytes)
			if !groupExp {
				for index, file := range De.Records {
					if file.Stu.StudentNum == StudentNum || file.Stu.StudentName == StudentName {
						if file.Stu.StudentNum != StudentNum || file.Stu.StudentName != StudentName {
							Logger.LogW("detect wrong named file '%s':\n"+
								">> DataBase: %s, %s\n"+
								">> FileName: %s, %s", fullName, file.Stu.StudentNum, file.Stu.StudentName,
								StudentNum, StudentName)
						}
						De.Records[index].LastUpload = Time
						De.Records[index].Uploaded = true
						De.Records[index].Hash = Hash
						break
					}
				}
			} else {
				De.Records[StudentGrp-1].LastUpload = Time
				De.Records[StudentGrp-1].Uploaded = true
				De.Records[StudentGrp-1].Hash = Hash
			}
		}
		DynamicExpInfo = append(DynamicExpInfo, De)
	}
	return nil
}

func updateStudent(stuNum string, stuName string, exp int, time time.Time, hash string) bool {
	if !studentNumNameCheck(stuNum, stuName) {
		return false
	}
	index := StudentNumMap[stuNum]
	DynamicExpInfo[exp-1].Records[index].LastUpload = time
	DynamicExpInfo[exp-1].Records[index].Uploaded = true
	DynamicExpInfo[exp-1].Records[index].Hash = hash
	return true
}

func updateGroup(groupNum int, exp int, time time.Time, hash string) bool {
	if groupNum > len(GroupStu) || groupNum < 1 {
		return false
	}
	DynamicExpInfo[exp-1].Records[groupNum-1].LastUpload = time
	DynamicExpInfo[exp-1].Records[groupNum-1].Uploaded = true
	DynamicExpInfo[exp-1].Records[groupNum-1].Hash = hash
	return true
}
