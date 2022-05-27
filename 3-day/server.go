package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

func startMyGame(temp *Dealer, gameResult chan string) {
	wait := new(sync.WaitGroup)
	wait.Add(1)
	// mutex := sync.RWMutex{}

	go func() {
		defer wait.Done()
		for {
			select {
			case <-time.After(10 * time.Second):
				fmt.Println("timeout 10")
				if len(temp.players) < 2 {
					fmt.Println("No game, Player Number is : ", len(temp.players))

				} else if len(temp.players) > 6 {
					fmt.Println("No game, Player Number is : ", len(temp.players))
					s1 := rand.NewSource(time.Now().UnixNano())
					rand := rand.New(s1)
					delPlayerIndex := rand.Intn(len(temp.players))

					temp.deletePlayer(temp.players[delPlayerIndex].name)

				} else {
					temp.makeCard()
					temp.cardShuffle()
					temp.startGame()
					gameResultData := temp.printPlayerStatus()
					for range temp.players {
						gameResult <- gameResultData
					}
				}
			}
		}
	}()

	wait.Wait()

}
func main() {

	fmt.Println("Start cardGame!")

	dealer := Dealer{name: "dealer"}
	dealer.makeCard()

	networkServer(&dealer)

}

func networkServer(dealer *Dealer) {
	myListen, err := net.Listen("tcp", ":5000")

	gameResult := make(chan string, 5)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wait := new(sync.WaitGroup)
	wait.Add(1)
	go startMyGame(dealer, gameResult)

	for {
		connect, err := myListen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go ConnectHandler(connect, dealer, gameResult)
	}

	wait.Wait()
	defer func() {
		myListen.Close()

	}()
}

func ConnectHandler(connect net.Conn, dealer *Dealer, gameResult chan string) {
	recvBuf := make([]byte, 4096) // receive buffer: 4kB
	playerName := ""
	fmt.Println("Connect User!!")
	defer func() {
		connect.Close()
		fmt.Println(playerName + " ConnectHandler End")
	}()

	n, err := connect.Read(recvBuf)
	if err != nil {
		if io.EOF == err {
			fmt.Println("connection is closed from client; %v", connect.RemoteAddr().String())
			return
		}
		fmt.Println(err)
		return
	}

	if 0 < n {
		data := recvBuf[:n]
		if len(playerName) == 0 {
			playerName = string(data)
			playerName = connect.RemoteAddr().String()
			player := Player{name: playerName, age: 30}
			dealer.addPlayer(player)
		}
	}

	connectFlag := make(chan bool)
	//wait := new(sync.WaitGroup)
	//wait.Add(1)
	go func() {
		defer func() {
			//wait.Done()
			dealer.deletePlayer(playerName)
			connectFlag <- false
		}()
		_, err := connect.Read(recvBuf)
		if err != nil {
			if io.EOF == err {
				fmt.Println("connection is closed from client")
				return
			}
			fmt.Println(err)
			return
		}
	}()

	for {
		select {
		case result := <-gameResult:
			n, err = connect.Write([]byte(result))
			if err != nil {
				if io.EOF == err {
					fmt.Println(err)
				}
			}
		case flag := <-connectFlag:
			if flag == false {
				return
			}
		}
	}
	//wait.Wait()
}

type DealerInterface interface {
	makeCard()
	cardShuffle()
	addPlayer(Player)
	startGame()
	checkGame()
	printPlayerStatus()
}
type Dealer struct {
	card    []int
	name    string
	players []Player
	round   int
	draw    int
}

func (d *Dealer) makeCard() {
	d.card = make([]int, 20)
	for i := 0; i < 20; i++ {
		d.card[i] = (i + 1) % 10
		if d.card[i] == 0 {
			d.card[i] = 10
		}
	}
}

func (d *Dealer) cardShuffle() {

	myCard := make([]int, 20)

	for index := range myCard {
		s1 := rand.NewSource(time.Now().UnixNano())
		rand := rand.New(s1)
		randomNumber := rand.Intn(len(d.card))
		myCard[index] = d.card[randomNumber]
		d.card = append(d.card[:randomNumber], d.card[randomNumber+1:]...)
	}

	d.card = myCard
}

func (d *Dealer) addPlayer(p Player) {
	d.players = append(d.players, p)
	fmt.Println("add Player!! ", p.name)
}

func (d *Dealer) deletePlayer(playerName string) {
	fmt.Println("Delete Player : ", playerName)
	flagIndex := -1
	for index, val := range d.players {
		if val.name == playerName {
			flagIndex = index
			break
		}
	}
	fmt.Println(d.players[flagIndex].name, flagIndex)

	d.players = append(d.players[:flagIndex], d.players[flagIndex+1:]...)

}

func (d *Dealer) startGame() {
	for index := range d.players {
		d.players[index].mycard = append(d.players[index].mycard, d.card[(index*2):(index*2)+2])
	}
	d.checkGame()

}

func (d *Dealer) checkGame() {
	playerResult := make([]int, 0)
	// fmt.Println(d.players)
	for index := range d.players {
		cardIndex := len(d.players[index].mycard) - 1
		playerResult = append(playerResult, (int(d.players[index].mycard[cardIndex][0]+d.players[index].mycard[cardIndex][1]))%10)
	}
	max := -1
	draw_check := 0
	max_index := 0
	for i := range playerResult {
		temp := playerResult[i]
		if temp > max {
			max = temp
			max_index = i
		} else if max == temp {
			draw_check++
		} else {
			continue
		}
	}

	if draw_check == 0 {
		d.players[max_index].winHit++
	} else {
		d.draw++
	}
	d.round++
	// fmt.Println(playerResult)
}

func (d *Dealer) printPlayerStatus() string {
	result := ""
	for i := range d.players {
		fmt.Println("Player ", d.players[i].name, " WinHit : ", d.players[i].winHit)
		result += fmt.Sprintf("Player %s WinHit : %f\n", d.players[i].name, d.players[i].winHit)
	}
	fmt.Println("Round : ", d.round)
	fmt.Println("Draw Round: ", d.draw)
	result += fmt.Sprintf("Round : %d Draw Round: %d\n", d.round, d.draw)
	return result

}

type Player struct {
	name      string
	age       int
	mycard    [][]int
	winHit    float32
	dropRound int
}

func (p *Player) receiveCard(receive_card []int) {
	p.mycard = append(p.mycard, receive_card)
}
