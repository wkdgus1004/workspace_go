package main

import (
	"fmt"
	"math/rand"
)

func main() {

	var winRate [3]int
	winRate[0] = 0
	winRate[1] = 0
	winRate[2] = 0
	for i := 0; i < 100; i++ {
		//make Card Set
		var card [20]int
		for i := 0; i < 20; i++ {
			card[i] = (i + 1) % 10
			if card[i] == 0 {
				card[i] = 10
			}
		}

		//shuffle Card Set
		var myCard [20]int
		fmt.Println(card)
		var index = 0
		for {

			// Check Card
			flag := 0
			for _, temp := range card {
				flag = flag + temp
			}
			if flag == 0 || index == 20 {
				break
			}

			// s1 := rand.NewSource(time.Now().UnixNano())
			// rand := rand.New(s1)
			randomNumber := rand.Intn(20)
			if card[randomNumber] != 0 {
				myCard[index] = card[randomNumber]
				card[randomNumber] = 0
				index++
			} else {
				continue
			}

		}
		fmt.Println(myCard)

		var player_1, player_2 [2]int
		player_1[0] = myCard[0]
		player_2[0] = myCard[1]
		player_1[1] = myCard[2]
		player_2[1] = myCard[3]

		// compatition
		player_1_result := (player_1[0] + player_1[1]) % 10
		player_2_result := (player_2[0] + player_2[1]) % 10
		winner := 0
		if player_1_result > player_2_result {
			winner = 1
		} else if player_1_result == player_2_result {
			winner = 0
		} else {
			winner = 2
		}

		fmt.Println(player_1_result, player_2_result, winner)
		winRate[winner]++
	}

	fmt.Println(winRate)

	fmt.Printf("Player1 Win Rate : %d, %0.2f %%\n", winRate[1], float32(winRate[1]))
	fmt.Printf("Player2 Win Rate : %d, %0.2f %%\n", winRate[2], float32(winRate[2]))
	fmt.Printf("Draw Rate : %d, %0.2f %%\n", winRate[0], float32(winRate[0]))

}
