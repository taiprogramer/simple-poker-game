import 'package:flutter/material.dart';
import 'package:simple_poker_game/models/room.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';
import 'package:simple_poker_game/services/room/room_service.dart';
import 'dart:math' as math;

class RoomTexasHoldemPage extends StatefulWidget {
  static const String routeName = '/roomTexasHoldem';
  const RoomTexasHoldemPage({Key? key}) : super(key: key);

  @override
  State<RoomTexasHoldemPage> createState() => _RoomTexasHoldemPageState();
}

class _RoomTexasHoldemPageState extends State<RoomTexasHoldemPage> {
  Widget _playerInSlot({int slot = -1}) {
    final index = slot - 1; // because slot count from 1
    // current sign in user
    if (slot == 0) {
      return _PlayerCircle();
    }
    // slot is out of range
    if (slot > room.users.length) {
      return const Text("");
    }
    // skip current sign in user
    if (room.users.elementAt(index).id == AppLocalStorage.getItem('user_id')) {
      return const Text("");
    }

    return _PlayerCircle();
  }

  Room room = Room();

  Future<void> _fetchData() async {
    final roomData =
        await RoomService.getRoom(roomID: AppLocalStorage.getItem("room_id"));
    setState(() {
      room = roomData;
    });
  }

  @override
  void initState() {
    super.initState();
    _fetchData();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Simple Poker Game'),
      ),
      body: Column(
        children: [
          Container(
            padding: const EdgeInsets.all(20.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: const [
                Text(
                  'Pot: ',
                  style: TextStyle(fontSize: 24),
                ),
                Text(
                  '0',
                  style: TextStyle(fontSize: 24),
                ),
              ],
            ),
          ),
          SizedBox(
            height: 350,
            child: Container(
              decoration: const BoxDecoration(
                  image: DecorationImage(
                      image: AssetImage('assets/images/poker_table.jpg'))),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _playerInSlot(slot: 3),
                      Container(
                        child: _playerInSlot(slot: 1),
                        margin: const EdgeInsets.only(bottom: 40.0),
                      ),
                      Container(
                        child: _playerInSlot(slot: 2),
                        margin: const EdgeInsets.only(bottom: 40.0),
                      ),
                      _playerInSlot(slot: 4),
                    ],
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _playerInSlot(slot: 5),
                      _playerInSlot(slot: 6),
                    ],
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _playerInSlot(slot: 7),
                      Container(
                        margin: const EdgeInsets.only(top: 40.0),
                        child: _playerInSlot(slot: 0), // current sign in user
                      ),
                      _playerInSlot(slot: 8),
                    ],
                  )
                ],
              ),
            ),
          ),
          Container(
            padding: const EdgeInsets.all(10.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                ElevatedButton(onPressed: () {}, child: const Text('Fold')),
                ElevatedButton(onPressed: () {}, child: const Text('Call')),
                ElevatedButton(onPressed: () {}, child: const Text('Raise')),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(10.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                ElevatedButton(onPressed: () {}, child: const Text('Ready')),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(10.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                ElevatedButton(onPressed: () {}, child: const Text('Start')),
                ElevatedButton(onPressed: () {}, child: const Text('Delete')),
                ElevatedButton(onPressed: () {}, child: const Text('Delegate')),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

class _PlayerCircle extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
        width: 55,
        child: Stack(
          children: [
            Column(
              children: [
                Container(
                  width: 10,
                  height: 10,
                  decoration: const BoxDecoration(
                      image: DecorationImage(
                          image: AssetImage('assets/images/active_tick.png'))),
                ),
                Container(
                  alignment: Alignment.center,
                  height: 50,
                  width: 50,
                  decoration: BoxDecoration(
                      color: Colors.blue,
                      borderRadius: BorderRadius.circular(100)),
                  child: const Text(
                    'G',
                    style: TextStyle(fontWeight: FontWeight.bold),
                  ),
                ),
                const Text(
                  '\$ 999',
                  style: TextStyle(
                      color: Colors.red, backgroundColor: Colors.yellow),
                )
              ],
            ),
            Positioned(
                left: -10,
                child: Row(
                  children: [
                    Transform.rotate(
                      angle: -math.pi / 8,
                      child: Container(
                          width: 35,
                          height: 55,
                          decoration: const BoxDecoration(
                              image: DecorationImage(
                                  image: AssetImage(
                                      'assets/images/deck_of_cards/CLUB-1.png')))),
                    ),
                    Transform.rotate(
                      angle: math.pi / 8,
                      child: Container(
                          width: 35,
                          height: 55,
                          decoration: const BoxDecoration(
                              image: DecorationImage(
                                  image: AssetImage(
                                      'assets/images/deck_of_cards/SPADE-11-JACK.png')))),
                    )
                  ],
                )),
          ],
        ));
  }
}
