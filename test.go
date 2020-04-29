package ini

const testData string = `
GlobalKey=GlobalValue

[SectionName]
PropKeyString=PropValString
PropKeyInt=-1
PropKeyUint=1
PropKeyBool=true
PropKeyFloat=1.2
PropKeySlice=PropValSlice1
PropKeySlice=PropValSlice2
PropKeyMap[MapKey1]=PropValMap1
PropKeyMap[MapKey2]=PropValMap2

[SectionName2]
PropKeyString=PropValString
PropKeyInt=-1
PropKeyUint=1
PropKeyBool=true
PropKeyFloat=1.2
PropKeySlice=PropValSlice1
PropKeySlice=PropValSlice2
PropKeyMap[MapKey1]=PropValMap1
PropKeyMap[MapKey2]=PropValMap2

[SectionSlice1]
PropKeyString=PropValString

[SectionSlice1]
PropKeyString=PropValString`

type testDataSectionName struct {
	PropKeyString string
	PropKeyInt    int
	PropKeyUint   uint
	PropKeyBool   bool
	PropKeyFloat  float64
	PropKeySlice  []string
	PropKeyMap    map[string]string
}

type testDataSectionSlice struct {
	PropKeyString string
}

type testDataStruct struct {
	GlobalKey     string
	SectionName   testDataSectionName
	SectionName2  testDataSectionName
	SectionSlice1 []testDataSectionSlice
}

var want testDataStruct = testDataStruct{
	GlobalKey: "GlobalValue",
	SectionName: testDataSectionName{
		PropKeyString: "PropValString",
		PropKeyInt:    -1,
		PropKeyUint:   1,
		PropKeyBool:   true,
		PropKeyFloat:  1.2,
		PropKeySlice:  []string{"PropValSlice1", "PropValSlice2"},
		PropKeyMap: map[string]string{
			"MapKey1": "PropValMap1",
			"MapKey2": "PropValMap2",
		},
	},
	SectionName2: testDataSectionName{
		PropKeyString: "PropValString",
		PropKeyInt:    -1,
		PropKeyUint:   1,
		PropKeyBool:   true,
		PropKeyFloat:  1.2,
		PropKeySlice:  []string{"PropValSlice1", "PropValSlice2"},
		PropKeyMap: map[string]string{
			"MapKey1": "PropValMap1",
			"MapKey2": "PropValMap2",
		},
	},
	SectionSlice1: []testDataSectionSlice{
		{
			PropKeyString: "PropValString",
		},
		{
			PropKeyString: "PropValString",
		},
	},
}
