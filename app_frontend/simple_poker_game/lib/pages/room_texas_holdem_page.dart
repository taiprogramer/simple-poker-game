import 'package:flutter/material.dart';
import 'package:simple_poker_game/models/room.dart';
import 'package:simple_poker_game/models/table.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';
import 'package:simple_poker_game/services/room/room_service.dart';
import 'package:simple_poker_game/services/socket/socket.dart';
import 'dart:math' as math;

import 'package:simple_poker_game/services/table/table_service.dart';

class RoomTexasHoldemPage extends StatefulWidget {
  static const String routeName = '/roomTexasHoldem';
  const RoomTexasHoldemPage({Key? key}) : super(key: key);

  @override
  State<RoomTexasHoldemPage> createState() => _RoomTexasHoldemPageState();
}

class _RoomTexasHoldemPageState extends State<RoomTexasHoldemPage> {
  Room room = Room();
  bool ready = false;
  int tableID = 0;
  PokerTable table = PokerTable();

  final int userID = AppLocalStorage.getItem("user_id");
  final int roomID = AppLocalStorage.getItem("room_id");

  late SocketInstance socketInstance;

  Future<void> _refreshRoomState() async {
    final roomData =
        await RoomService.getRoom(roomID: AppLocalStorage.getItem('room_id'));
    late final bool readyStatus;
    for (final user in roomData.users) {
      if (user.id == userID) {
        readyStatus = user.ready;
      }
    }
    setState(() {
      room = roomData;
      ready = readyStatus;
    });
  }

  Future<void> _refreshTableState(int tableID) async {
    final tableData = await TableService.getTable(tableID, userID);
    setState(() {
      table = tableData;
    });
  }

  void _socketListener(String msg) async {
    if (msg == "new user join room" || msg == "room status was changed") {
      _refreshRoomState();
    }

    if (msg.startsWith("table=")) {
      final tableIDStr = msg.substring(msg.indexOf("=") + 1);
      tableID = int.parse(tableIDStr);
      _refreshTableState(tableID);
    }

    if (msg == "the game is started") {
      _refreshRoomState();
    }
  }

  void _connectWebSocket() {
    final userID = AppLocalStorage.getItem("user_id");
    final roomID = AppLocalStorage.getItem('room_id');
    socketInstance = SocketInstance(userID: userID, roomID: roomID);
    socketInstance.connect();
    socketInstance.listen(_socketListener);
  }

  void _startTheGame() {
    socketInstance.send("start");
  }

  String _buildImageUrl(int number, int suit) {
    final List<String> suits = ['DIAMOND', 'HEART', 'CLUB', 'SPADE'];
    return 'assets/images/deck_of_cards/${suits.elementAt(suit)}-$number.png';
  }

  Widget _playerInSlot({int slot = -1}) {
    final index = slot - 1; // because slot count from 1
    final userID = AppLocalStorage.getItem('user_id');
    String card1ImageUrl = '';
    String card2ImageUrl = '';
    // current sign in user
    if (slot == 0) {
      bool ready = false;
      for (final user in room.users) {
        if (user.id == userID) {
          ready = user.ready;
          if (room.playing) {
            final card1 = table.ownCards[0];
            final card2 = table.ownCards[1];
            card1ImageUrl = _buildImageUrl(card1.number, card1.suit);
            card2ImageUrl = _buildImageUrl(card2.number, card2.suit);
          }
          break;
        }
      }
      return _PlayerCircle(
        ready: ready,
        card1ImageUrl: card1ImageUrl,
        card2ImageUrl: card2ImageUrl,
      );
    }
    // slot is out of range
    if (slot > room.users.length) {
      return const Text("");
    }
    // skip current sign in user
    if (room.users.elementAt(index).id == userID) {
      return const Text("");
    }

    final ready = room.users.elementAt(index).ready;
    return _PlayerCircle(
      ready: ready,
      card1ImageUrl: card1ImageUrl,
      card2ImageUrl: card2ImageUrl,
    );
  }

  @override
  void initState() {
    super.initState();
    _refreshRoomState();
    _connectWebSocket();
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
        onWillPop: () async {
          socketInstance.disconnect();
          return true;
        },
        child: Scaffold(
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
                            child:
                                _playerInSlot(slot: 0), // current sign in user
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
                    ElevatedButton(
                        onPressed: () {}, child: const Text('Raise')),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.all(10.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    ElevatedButton(
                        onPressed: () async {
                          await RoomService.updateReadyStatus(
                              roomID: room.id, ready: !ready, userID: userID);
                          socketInstance.send("ready");
                        },
                        child: Text(ready ? 'Cancel' : 'Ready')),
                  ],
                ),
              ),
              Container(
                padding: const EdgeInsets.all(10.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    ElevatedButton(
                        onPressed: () {
                          _startTheGame();
                        },
                        child: const Text('Start')),
                    ElevatedButton(
                        onPressed: () {}, child: const Text('Delete')),
                    ElevatedButton(
                        onPressed: () {}, child: const Text('Delegate')),
                  ],
                ),
              ),
            ],
          ),
        ));
  }
}

class _PlayerCircle extends StatelessWidget {
  final String shortName;
  final int money;
  final String card1ImageUrl;
  final String card2ImageUrl;
  final bool active;
  final bool ready;

  const _PlayerCircle(
      {Key? key,
      this.shortName = 'G',
      this.money = 0,
      this.card1ImageUrl = '',
      this.card2ImageUrl = '',
      this.active = false,
      this.ready = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
        width: 55,
        child: Stack(
          children: [
            Column(
              children: [
                active
                    ? Container(
                        width: 10,
                        height: 10,
                        decoration: const BoxDecoration(
                            image: DecorationImage(
                                image: AssetImage(
                                    'assets/images/active_tick.png'))))
                    : const Text(''),
                Container(
                  alignment: Alignment.center,
                  height: 50,
                  width: 50,
                  decoration: BoxDecoration(
                      color: ready ? Colors.blue : Colors.grey,
                      borderRadius: BorderRadius.circular(100)),
                  child: Text(
                    shortName,
                    style: const TextStyle(fontWeight: FontWeight.bold),
                  ),
                ),
                Text(
                  '\$ $money',
                  style: const TextStyle(
                      color: Colors.red, backgroundColor: Colors.yellow),
                )
              ],
            ),
            Positioned(
                left: -10,
                child: Row(
                  children: [
                    card1ImageUrl != ''
                        ? Transform.rotate(
                            angle: -math.pi / 8,
                            child: Container(
                                width: 35,
                                height: 55,
                                decoration: BoxDecoration(
                                    image: DecorationImage(
                                        image: AssetImage(card1ImageUrl)))),
                          )
                        : const Text(''),
                    card2ImageUrl != ''
                        ? Transform.rotate(
                            angle: math.pi / 8,
                            child: Container(
                                width: 35,
                                height: 55,
                                decoration: BoxDecoration(
                                    image: DecorationImage(
                                        image: AssetImage(card2ImageUrl)))),
                          )
                        : const Text('')
                  ],
                )),
          ],
        ));
  }
}
