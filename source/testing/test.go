package main

import (
	//"fmt"
	. "source/config"
	//"source/localElevator/fsm"
	//"source/network/peers"
	//"source/primary/assigner"
)

func main(){
	Ids:=[]string{"Ids[0]","Ids[1]"}
	/* for index,Id:=range(Ids){
		fmt.Printf("ID nr %d = %s\n",index,Id)
	} */
	
	el:=make(map[string]Elevator)
												 // HALLUP HALLDWN  CAB
	el[Ids[0]]=Elevator{Floor:2,Orders:[4][3]bool{{false, true, false}, //FLOOR 4
													{false, false, true}, //FLOOR 3
													{true, true, false}, //FLOOR 2
													{false, false, true}, //FLOOR 1
													},
				PrevDirection:UP}
	el[Ids[1]]=Elevator{Floor:3,Orders:[4][3]bool{{false, false, false},
													{false, true, false},
													{true, true, false},
													{false, false, false},
													},
				PrevDirection:UP}
	el["ok"]=Elevator{Floor:2,Orders:[4][3]bool{	{false, false, false},
													{false, false, false},
													{false, false, false},
													{false, false, false}, 
												},
				PrevDirection:UP}
	el["ok2"]=Elevator{Floor:3,Orders:[4][3]bool{	{false, false, false},
													{false, false, false},
													{false, false, false},
													{false, false, false},
												},
				PrevDirection:UP}			

	//p := peers.PeerUpdate{Ids,"hei",[]string{"hei","hei"}}
	// activeIds := []string{Ids[0],Ids[1],"ok2"}//p.Peers

	// fmt.Println("Time until pickup for el[Ids[0]]: ",fsm.TimeUntilPickup(el[Ids[0]],Order{0,1}))
	// fmt.Println("Time until pickup for el[Ids[1]]: ",fsm.TimeUntilPickup(el[Ids[1]],Order{0,1}))
	// fmt.Println("Time until pickup for el['ok']: ",fsm.TimeUntilPickup(el["ok"],Order{0,1}))
	// fmt.Println("Time until pickup for el['ok2']: ",fsm.TimeUntilPickup(el["ok2"],Order{0,1}))
	
	// fmt.Println("Best elevator: ",assigner.ChooseElevator(el,activeIds,Order{0,1}))
}