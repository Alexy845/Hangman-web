package hangmanweb

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func InitGlobalValue(w http.ResponseWriter, r *http.Request, globaldata GlobalInfo) GlobalInfo {

	globaldata.Username = GetCookieAccount(r)

	globaldata.Status = GetCookieStatus(r)

	if globaldata.Status == "login" && UserExist(globaldata.Username) {
		globaldata.UserLevel, globaldata.UserXpAv = AtoiWithoutErr(GetUserInfo(GetCookieAccount(r))[5]), float64(AtoiWithoutErr(GetUserInfo(GetCookieAccount(r))[6])/AtoiWithoutErr(GetUserInfo(GetCookieAccount(r))[5]))
	}

	if !UserExist(globaldata.Username) {
		SetCookieAccount(w, "", "logout")
	}

	globalDatabase, err := os.OpenFile("./server/database/global.csv", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	defer globalDatabase.Close()

	csvReaderGlobalDB := csv.NewReader(globalDatabase)
	getDataGlobalDB, err := csvReaderGlobalDB.ReadAll()

	if err != nil {
		fmt.Println(err)
	}

	if len(getDataGlobalDB) != 0 {
		globaldata.DeadSanta = AtoiWithoutErr(getDataGlobalDB[0][0])
		globaldata.SaveSanta = AtoiWithoutErr(getDataGlobalDB[0][1])
	} else {

		csvWriterGlobalDB := csv.NewWriter(globalDatabase)

		newData := []string{"0", "0"}
		err = csvWriterGlobalDB.Write(newData)
		if err != nil {
			fmt.Println(err)
		}
		defer csvWriterGlobalDB.Flush()

		globaldata.DeadSanta = 0
		globaldata.SaveSanta = 0
	}

	if globaldata.SaveSanta+globaldata.DeadSanta != 0 {
		globaldata.Ratio = globaldata.SaveSanta * 100 / (globaldata.SaveSanta + globaldata.DeadSanta)
	} else {
		globaldata.Ratio = 50
	}

	globaldata.Total = globaldata.DeadSanta + globaldata.SaveSanta
	return globaldata
}

func UpdateGlobalValue(w http.ResponseWriter, r *http.Request, save bool, globaldata GlobalInfo) GlobalInfo {

	globaldata.Username = GetCookieAccount(r)

	globaldata.Status = GetCookieStatus(r)

	if globaldata.Status == "login" && UserExist(globaldata.Username) {
		globaldata.UserLevel, globaldata.UserXpAv = AtoiWithoutErr(GetUserInfo(GetCookieAccount(r))[5]), float64(AtoiWithoutErr(GetUserInfo(GetCookieAccount(r))[6])/AtoiWithoutErr(GetUserInfo(GetCookieAccount(r))[5]))
	}

	if !UserExist(globaldata.Username) {
		SetCookieAccount(w, "", "logout")
	}

	globalDatabase, err := os.OpenFile("./server/database/global.csv", os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	csvReaderGlobalDB := csv.NewReader(globalDatabase)
	getDataGlobalDB, err := csvReaderGlobalDB.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	if len(getDataGlobalDB) != 0 {
		globaldata.DeadSanta = AtoiWithoutErr(getDataGlobalDB[0][0])
		globaldata.SaveSanta = AtoiWithoutErr(getDataGlobalDB[0][1])
	} else {
		globaldata.DeadSanta = 0
		globaldata.SaveSanta = 0
	}
	globalDatabase.Close()

	globalDatabase, err = os.OpenFile("./server/database/global.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	csvWriterGlobalDB := csv.NewWriter(globalDatabase)

	if save {
		globaldata.SaveSanta++
	} else {
		globaldata.DeadSanta++
	}

	newData := []string{strconv.Itoa(globaldata.DeadSanta), strconv.Itoa(globaldata.SaveSanta)}
	err = csvWriterGlobalDB.Write(newData)
	if err != nil {
		fmt.Println(err)
	}
	defer csvWriterGlobalDB.Flush()

	if globaldata.SaveSanta+globaldata.DeadSanta != 0 {
		globaldata.Ratio = globaldata.SaveSanta * 100 / (globaldata.SaveSanta + globaldata.DeadSanta)
	} else {
		globaldata.Ratio = 50
	}

	globaldata.Total = globaldata.DeadSanta + globaldata.SaveSanta

	return globaldata

}

func UpdateUserValue(win bool, w http.ResponseWriter, r *http.Request, sbUsersList ScoreboardData, gameLaunch map[string]Hangman) ScoreboardData {
	userDatabase, err := os.OpenFile("./server/database/users.csv", os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	csvReaderUsersDB := csv.NewReader(userDatabase)
	getDataUsersDB, err := csvReaderUsersDB.ReadAll()

	sbUsersList.UsersList = []User{}

	if len(getDataUsersDB) != 0 {
		for ligne, userIngetData := range getDataUsersDB {
			if userIngetData[0] == gameLaunch[CookieSession(w, r, gameLaunch)].PlayerName {
				multi := 0
				levelGain := 0
				switch gameLaunch[CookieSession(w, r, gameLaunch)].Mode {
				case "easy":
					multi = 1
				case "medium":
					multi = 2
				case "hard":
					multi = 3
				}
				if win {
					xpGain := multi*gameLaunch[CookieSession(w, r, gameLaunch)].Attempts + multi
					if xpGain+AtoiWithoutErr(userIngetData[6]) >= AtoiWithoutErr(userIngetData[5])*100 {
						levelGain++
					}
					// 			attempts * multi + multi
					sbUsersList.UsersList = append(sbUsersList.UsersList, User{userIngetData[0], AtoiWithoutErr(userIngetData[2]) + 1, AtoiWithoutErr(userIngetData[3]), AtoiWithoutErr(userIngetData[4]) + 1, AtoiWithoutErr(userIngetData[5]) + levelGain, AtoiWithoutErr(userIngetData[6]) + xpGain - (((AtoiWithoutErr(userIngetData[5])) * 100) * levelGain)})
					getDataUsersDB[ligne][2] = strconv.Itoa(AtoiWithoutErr(getDataUsersDB[ligne][2]) + 1)
					getDataUsersDB[ligne][4] = strconv.Itoa(AtoiWithoutErr(getDataUsersDB[ligne][4]) + 1)
					getDataUsersDB[ligne][6] = strconv.Itoa(AtoiWithoutErr(userIngetData[6]) + xpGain - (((AtoiWithoutErr(userIngetData[5])) * 100) * levelGain))
					getDataUsersDB[ligne][5] = strconv.Itoa(AtoiWithoutErr(userIngetData[5]) + levelGain)
				} else {
					xpGain := multi
					if xpGain+AtoiWithoutErr(userIngetData[6]) >= AtoiWithoutErr(userIngetData[5])*100 {
						levelGain++
					}
					// 			+ multi
					sbUsersList.UsersList = append(sbUsersList.UsersList, User{userIngetData[0], AtoiWithoutErr(userIngetData[2]) + 1, AtoiWithoutErr(userIngetData[3]), AtoiWithoutErr(userIngetData[4]) + 1, AtoiWithoutErr(userIngetData[5]) + levelGain, AtoiWithoutErr(userIngetData[6]) + xpGain - (((AtoiWithoutErr(userIngetData[5])) * 100) * levelGain)})
					getDataUsersDB[ligne][3] = strconv.Itoa(AtoiWithoutErr(getDataUsersDB[ligne][3]) + 1)
					getDataUsersDB[ligne][4] = strconv.Itoa(AtoiWithoutErr(getDataUsersDB[ligne][4]) + 1)
					getDataUsersDB[ligne][6] = strconv.Itoa(AtoiWithoutErr(userIngetData[6]) + xpGain - (((AtoiWithoutErr(userIngetData[5])) * 100) * levelGain))
					getDataUsersDB[ligne][5] = strconv.Itoa(AtoiWithoutErr(userIngetData[5]) + levelGain)
					//5 = level
					//6 = exp		+ multi
				}
			} else {
				sbUsersList.UsersList = append(sbUsersList.UsersList, User{userIngetData[0], AtoiWithoutErr(userIngetData[2]), AtoiWithoutErr(userIngetData[3]), AtoiWithoutErr(userIngetData[4]), AtoiWithoutErr(userIngetData[5]), AtoiWithoutErr(userIngetData[6])})
			}
		}
	}
	userDatabase.Close()

	userDatabase, err = os.OpenFile("./server/database/users.csv", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	csvWriterUsersDB := csv.NewWriter(userDatabase)

	err = csvWriterUsersDB.WriteAll(getDataUsersDB)
	if err != nil {
		fmt.Println(err)
	}
	defer csvWriterUsersDB.Flush()

	return sbUsersList
}

func UserExist(username string) bool {
	usersDatabase, err := os.OpenFile("./server/database/users.csv", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	defer usersDatabase.Close()

	csvReaderUserslDB := csv.NewReader(usersDatabase)
	getDataUsersDB, err := csvReaderUserslDB.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range getDataUsersDB {
		if user[0] == username {
			return true
		}
	}
	return false

}

func GetUserInfo(username string) []string {
	usersDatabase, err := os.OpenFile("./server/database/users.csv", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}

	defer usersDatabase.Close()

	csvReaderUserslDB := csv.NewReader(usersDatabase)
	getDataUsersDB, err := csvReaderUserslDB.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, user := range getDataUsersDB {
		if user[0] == username {
			return user
		}
	}
	return []string{}

}
