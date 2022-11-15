import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_dialogs/flutter_dialogs.dart';
import 'package:simple_poker_game/models/room.dart';
import 'package:simple_poker_game/models/table.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';
import 'package:simple_poker_game/services/room/room_service.dart';
import 'package:simple_poker_game/services/socket/socket.dart';
import 'dart:math' as math;

import 'package:simple_poker_game/services/table/table_service.dart';

class Chat {
  final int userID;
  final String content;
  final bool seen;

  Chat({
    this.userID = 0,
    this.content = '',
    this.seen = false,
  });

  Map<String, dynamic> toJson() =>
      {'user_id': userID, 'content': content, 'seen': seen};

  factory Chat.fromMap(Map data) {
    return Chat(
      userID: data['user_id'],
      content: data['content'],
      seen: data['seen'],
    );
  }
}

class RoomTexasHoldemPage extends StatefulWidget {
  static const String routeName = '/roomTexasHoldem';
  const RoomTexasHoldemPage({Key? key}) : super(key: key);

  @override
  State<RoomTexasHoldemPage> createState() => _RoomTexasHoldemPageState();
}

class _RoomTexasHoldemPageState extends State<RoomTexasHoldemPage> {
  Room room = Room();
  bool ready = false;
  bool endGame = false;
  int tableID = 0;
  PokerTable table = PokerTable(currentTurn: UserTurn());
  List<Chat> chatMessages = List.empty(growable: true);

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

  Future<void> _refreshTableAndRoomState(int tableID, int roomID) async {
    final roomData = await RoomService.getRoom(roomID: roomID);
    final tableData = await TableService.getTable(tableID, userID);

    setState(() {
      room = roomData;
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
      _refreshTableAndRoomState(tableID, AppLocalStorage.getItem('room_id'));
    }

    if (msg == "the game has ended") {
      endGame = true;
      _refreshTableAndRoomState(tableID, AppLocalStorage.getItem('room_id'));
    }

    if (msg.startsWith("broadcast=")) {
      final msgPattern = msg.split("=").elementAt(1);
      final userID = msgPattern.split("\$").elementAt(0);
      final content = msgPattern.split("\$").elementAt(1);
      setState(() {
        chatMessages.add(
            Chat(userID: int.parse(userID), content: content, seen: false));
      });
      Timer timer = Timer(const Duration(seconds: 3), () {
        setState(() {
          chatMessages.removeAt(chatMessages.length - 1);
          chatMessages.add(
              Chat(userID: int.parse(userID), content: content, seen: true));
        });
      });
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

  void _performAction(String actionName, int amount) async {
    await TableService.performAction(table.id, userID, actionName, amount);
    socketInstance.send('has performed action');
  }

  String _buildImageUrl(int number, int suit) {
    final List<String> suits = ['DIAMOND', 'HEART', 'CLUB', 'SPADE'];
    return 'assets/images/deck_of_cards/${suits.elementAt(suit)}-$number.png';
  }

  String _getBackCardImageUrl() {
    return 'assets/images/back_card.png';
  }

  int _getAmount(PokerTable table, int userID) {
    for (int i = 0; i < table.players.length; i++) {
      final player = table.players.elementAt(i);
      if (player.id == userID) {
        return player.availableMoney;
      }
    }
    return 0;
  }

  String _getLatestChatUnseen(int userID) {
    String msg = '';

    for (var i = chatMessages.length - 1; i >= 0; i--) {
      Chat chat = chatMessages.elementAt(i);
      if (chat.userID == userID) {
        if (!chat.seen) {
          msg = chat.content;
        }
        break;
      }
    }

    return msg;
  }

  Widget _playerInSlot({int slot = -1}) {
    // because slot count from 1 except current sign-in user slot is 0.
    final index = slot == 0 ? 0 : slot - 1;
    final userID = AppLocalStorage.getItem('user_id');
    String card1ImageUrl = room.playing ? _getBackCardImageUrl() : '';
    String card2ImageUrl = room.playing ? _getBackCardImageUrl() : '';
    int amount = 0;
    bool active = false;
    String msg = '';
    // current sign in user
    if (slot == 0) {
      bool ready = false;
      String shortName = '';
      for (final user in room.users) {
        if (user.id == userID) {
          msg = _getLatestChatUnseen(userID);
          shortName =
              user.username.substring(user.username.length - 1).toUpperCase();
          ready = user.ready;
          if (room.playing) {
            final card1 = table.ownCards[0];
            final card2 = table.ownCards[1];
            card1ImageUrl = _buildImageUrl(card1.number, card1.suit);
            card2ImageUrl = _buildImageUrl(card2.number, card2.suit);
            active = table.currentTurn.userID == userID;
            amount = _getAmount(table, userID);
          }
          break;
        }
      }
      return _PlayerCircle(
        ready: ready,
        card1ImageUrl: card1ImageUrl,
        card2ImageUrl: card2ImageUrl,
        active: active,
        shortName: shortName,
        money: amount,
        msg: msg,
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
    final user = room.users.elementAt(index);
    final shortName =
        user.username.substring(user.username.length - 1).toUpperCase();
    if (room.playing) {
      active = table.currentTurn.userID == room.users[index].id;
      amount = _getAmount(table, user.id);
    }
    msg = _getLatestChatUnseen(user.id);
    if (endGame) {
      // if the game is end, there is no active tick
      active = false;
      // show 2 cards of other players as well if the game is ended
      for (int i = 0; i < table.results.length; i++) {
        final result = table.results.elementAt(i);
        if (result.userID == user.id) {
          final card1 = result.cards.elementAt(0);
          final card2 = result.cards.elementAt(1);
          card1ImageUrl = _buildImageUrl(card1.number, card1.suit);
          card2ImageUrl = _buildImageUrl(card2.number, card2.suit);
          break;
        }
      }
    }

    return _PlayerCircle(
      ready: ready,
      card1ImageUrl: card1ImageUrl,
      card2ImageUrl: card2ImageUrl,
      active: active,
      shortName: shortName,
      money: amount,
      msg: msg,
    );
  }

  List<Widget> _generateCommonCardWidgets() {
    List<Widget> commonCardWidgets = List.empty(growable: true);
    for (int i = 0; i < table.commonCards.length; i++) {
      final card = table.commonCards.elementAt(i);
      final cardImageUrl = _buildImageUrl(card.number, card.suit);
      commonCardWidgets.add(Container(
        height: 80.0,
        width: 30.0,
        decoration: BoxDecoration(
            image: DecorationImage(image: AssetImage(cardImageUrl))),
      ));
    }
    return commonCardWidgets;
  }

  Widget _commonCards() {
    return Row(children: _generateCommonCardWidgets());
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
          body: Column(
            children: [
              Container(
                padding: const EdgeInsets.all(20.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Row(
                      children: [
                        const Text(
                          'Round: ',
                          style: TextStyle(fontSize: 24),
                        ),
                        Text(
                          table.round.toString(),
                          style: const TextStyle(fontSize: 24),
                        ),
                      ],
                    ),
                    Row(
                      children: [
                        const Text(
                          'Pot: ',
                          style: TextStyle(fontSize: 24),
                        ),
                        Text(
                          table.pot.toString(),
                          style: const TextStyle(fontSize: 24),
                        ),
                      ],
                    )
                  ],
                ),
              ),
              Stack(children: [
                SizedBox(
                  height: 300,
                  child: Container(
                    decoration: const BoxDecoration(
                        image: DecorationImage(
                            image:
                                AssetImage('assets/images/poker_table.jpg'))),
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
                            _commonCards(),
                            _playerInSlot(slot: 6),
                          ],
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceAround,
                          children: [
                            _playerInSlot(slot: 7),
                            Container(
                              margin: const EdgeInsets.only(top: 40.0),
                              child: _playerInSlot(
                                  slot: 0), // current sign in user
                            ),
                            _playerInSlot(slot: 8),
                          ],
                        )
                      ],
                    ),
                  ),
                ),
                Positioned(
                    right: 5.0,
                    top: 5.0,
                    child: IconButton(
                      onPressed: () {
                        showPlatformDialog(
                            context: context,
                            builder: (context) {
                              var msg = '';
                              return BasicDialogAlert(
                                content: Row(
                                  children: [
                                    SizedBox(
                                      width: 300.0,
                                      child: TextFormField(
                                          onChanged: (String? value) {
                                            msg = value ?? '';
                                          },
                                          decoration: const InputDecoration(
                                              labelText: 'Message',
                                              icon: Icon(
                                                Icons.mail_outline,
                                                color: Colors.pinkAccent,
                                              ))),
                                    ),
                                    IconButton(
                                      icon: const Icon(Icons.send,
                                          color: Colors.pink),
                                      onPressed: () {
                                        socketInstance.send('broadcast=$msg');
                                        Navigator.pop(context);
                                      },
                                    ),
                                  ],
                                ),
                              );
                            });
                      },
                      icon: const Icon(Icons.message, color: Colors.pink),
                    )),
                Positioned(
                  bottom: 5.0,
                  right: 10.0,
                  child: Row(children: [
                    ElevatedButton(
                        onPressed: () {
                          _performAction("check", 0);
                        },
                        child: const Text('Check')),
                    ElevatedButton(
                        onPressed: () {
                          _performAction("fold", 0);
                        },
                        child: const Text('Fold')),
                    ElevatedButton(
                        onPressed: () {
                          // zero have no meaning here because call action
                          // doesn't need amount value.
                          _performAction('call', 0);
                        },
                        child: const Text('Call')),
                    ElevatedButton(
                        onPressed: () {
                          _performAction('raise', 100);
                        },
                        child: const Text('Raise')),
                  ]),
                ),
                Positioned(
                    left: 10.0,
                    bottom: 5.0,
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
                    )),
                Positioned(
                    right: 10.0,
                    bottom: 50.0,
                    child: ElevatedButton(
                        onPressed: () async {
                          await RoomService.updateReadyStatus(
                              roomID: room.id, ready: !ready, userID: userID);
                          socketInstance.send("ready");
                        },
                        child: Text(ready ? 'Cancel' : 'Ready'))),
              ]),
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
  final String msg;

  const _PlayerCircle(
      {Key? key,
      this.shortName = 'G',
      this.money = 0,
      this.card1ImageUrl = '',
      this.card2ImageUrl = '',
      this.active = false,
      this.ready = false,
      this.msg = ''})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
        width: 100.0,
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
            Positioned(
                child: Text(msg,
                    style: const TextStyle(
                        fontSize: 16.0,
                        backgroundColor: Colors.pinkAccent,
                        color: Colors.white))),
          ],
        ));
  }
}
