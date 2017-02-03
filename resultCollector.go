package main


import (
    "log"
    "io/ioutil"
    "bufio"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "os"
    "strings"
)

type ResultJson struct {
		Labels  string
    Timestamp   float32
    Metric string
    Official bool
		Value float32
		Test string
		Run_uri string
		Owner string
		Product_name string
		Unit string
		Sample_uri string
}

type Result struct {
    Timestamp   float32
    Metric string
    Official bool
		Value float32
		Test string
		Run_uri string
		Owner string
		Product_name string
		Unit string
		Sample_uri string
		Run_id int
		Labels  map[string]string
}


func collectResults(runCommand RunCommand)  {
  config := readConfig()
  rootDir:="/tmp/perfkitbenchmarker/runs/"
  files, _ := ioutil.ReadDir(rootDir)

  for _, f := range files {
   log.Printf(f.Name())
    if(f.IsDir()){
      jsonResults := readFile(rootDir+f.Name()+"/perfkitbenchmarker_results.json", rootDir+f.Name(), config.ResultFolder + string(runCommand.RunId) +"_" +f.Name())
      results := convertResults(jsonResults, runCommand.RunId)
      insertDB(results, config)
    }
  }
}

func moveFailed(runCommand RunCommand){
    config := readConfig()
    rootDir:="/tmp/perfkitbenchmarker/runs/"

    files, _ := ioutil.ReadDir(rootDir)
    for _, f := range files {
     log.Printf(f.Name())
      if(f.IsDir()){
          moveFolder(rootDir+f.Name(), config.ResultFolder +"failed/"+string(runCommand.RunId)+"_"+ f.Name())
      }
    }
}

func readFile(filePath string, srcPath string, dstPath string)([]ResultJson){

  jsonResults := []ResultJson{}

  file, err := os.Open(filePath)
  if err != nil {
    failOnError(err, "Couldn't result File")
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
      text:= scanner.Text()
      result := ResultJson{}
      err:= json.Unmarshal([]byte(text), &result)
      if(err != nil){
        failOnError(err, "Couldn't parse results")
      }
      jsonResults = append(jsonResults, result)
  }
  go moveFolder(srcPath, dstPath)

  return jsonResults
}

func moveFolder(srcFolder string, dstFolder string){
	if err:= os.Rename(srcFolder, dstFolder); err != nil {
    failOnError(err, "Couldn't move results folder")
	}
}

func convertResults(jsonResults []ResultJson, runId int)([]Result){
	results := []Result{}

	for _, jsonResult := range jsonResults {
		log.Printf("Converting: %s", jsonResult.Run_uri)
		result := Result{}

		result.Timestamp 		= jsonResult.Timestamp
		result.Metric				= jsonResult.Metric
		result.Official			= jsonResult.Official
		result.Value					= jsonResult.Value
		result.Test					= jsonResult.Test
		result.Run_uri				= jsonResult.Run_uri
		result.Owner					= jsonResult.Owner
		result.Product_name	= jsonResult.Product_name
		result.Unit					= jsonResult.Unit
		result.Sample_uri		= jsonResult.Sample_uri
		result.Run_id 				= runId
		result.Labels 				= extractLabels(jsonResult)

    if(result.Labels["num_parallel_copies"] != "") {
      result.Metric = result.Metric + " (num_parallel_copies=" + result.Labels["num_parallel_copies"] + ")"
    }

		results = append(results, result)

	}

	return results
}

func extractLabels(jsonResult ResultJson)(map[string]string){
	labels:= make(map[string]string)
	labelTupel := strings.Split(jsonResult.Labels, ",")
	for _, labelString := range labelTupel {
		labelString = strings.Replace(labelString, "|", "", -1)
		keyValue := strings.Split(labelString,":")
		labels[keyValue[0]] = keyValue[1]
	}
	return labels
}

func insertDB(results []Result, config Configuration){
   mongoDBDialInfo := &mgo.DialInfo{
  	Addrs:    []string{config.MongoHost},
  	Database: config.AuthDatabase,
  	Username: config.AuthUserName,
  	Password: config.AuthPassword,
  }

	session, err := mgo.DialWithInfo(mongoDBDialInfo)
	defer session.Close()
	if err != nil {
							panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB(config.Database).C(config.Collection)
	for _, result := range results {
		err = c.Insert(result)
		if err != nil {
				log.Fatalf("%s", err)
		}
	}
}
