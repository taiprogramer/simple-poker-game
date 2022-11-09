import 'package:flutter/material.dart';
import 'package:flutter_dialogs/flutter_dialogs.dart';
import 'package:simple_poker_game/models/user.dart';
import 'package:simple_poker_game/pages/room_texas_holdem_page.dart';
import 'package:simple_poker_game/services/auth/auth_service.dart';
import 'package:simple_poker_game/services/local_storage/local_storage.dart';

import '../models/room.dart';
import '../services/room/room_service.dart';

class TexasHoldemPage extends StatefulWidget {
  static const String routeName = '/texasHoldem';

  const TexasHoldemPage({Key? key, this.title = 'Simple Poker Game'})
      : super(key: key);

  final String title;

  @override
  State<TexasHoldemPage> createState() => _TexasHoldemPageState();
}

class _TexasHoldemPageState extends State<TexasHoldemPage> {
  late Future<List<Room>> rooms;
  late User user;
  String roomFilter = '';

  @override
  void initState() {
    super.initState();
    rooms = _fetchListRoom();
    _updateUserState();
  }

  Future<List<Room>> _fetchListRoom() {
    return RoomService.listRoom(0, 8);
  }

  void _updateUserState() async {
    int userID = AppLocalStorage.getItem('user_id');
    final userData = await AuthService.getUser(userID);
    setState(() {
      user = userData;
    });
  }

  void _updateListRoomState() {
    setState(() {
      rooms = _fetchListRoom();
    });
  }

  _onSliderChangedEndHandler(double value) async {
    await AppLocalStorage.setItem('amount', value);
  }

  _newRoom() async {
    try {
      double amount = AppLocalStorage.getItem('amount');
      int intAmount = amount.toInt();
      if (intAmount == 0) {
        return;
      }
      final room = await RoomService.newRoom(userID: user.id, money: intAmount);
      await AppLocalStorage.setItem('room_id', room.id);
      Navigator.pushNamed(context, RoomTexasHoldemPage.routeName);
    } catch (_) {}
  }

  _joinRoom() async {
    double amount = AppLocalStorage.getItem('amount');
    int intAmount = amount.toInt();
    if (intAmount == 0) {
      return;
    }
    final roomID = AppLocalStorage.getItem("room_id");
    final userID = AppLocalStorage.getItem('user_id');
    await RoomService.joinRoom(roomID: roomID, userID: userID);
    Navigator.pushNamed(context, RoomTexasHoldemPage.routeName);
  }

  _showAskingMoneyDialog(BuildContext context, Function runWhenGo) {
    showPlatformDialog(
      context: context,
      builder: (context) => BasicDialogAlert(
        title: const Text('How much money you want to bring to the room?'),
        content: SizedBox(
          height: 100,
          child: _MySlider(
            maxValue: user.money.toDouble(),
            onChangedEndHandler: _onSliderChangedEndHandler,
          ),
        ),
        actions: <Widget>[
          BasicDialogAction(
            title: const Text('Cancel'),
            onPressed: () async {
              await AppLocalStorage.setItem('amount', 0.0);
              Navigator.pop(context);
            },
          ),
          BasicDialogAction(
            title: const Text('Go'),
            onPressed: () async {
              Navigator.pop(context);
              runWhenGo();
              await AppLocalStorage.setItem('amount', 0.0);
            },
          ),
        ],
      ),
    );
  }

  _onTap(BuildContext context, int roomID) async {
    await AppLocalStorage.setItem("room_id", roomID);
    _showAskingMoneyDialog(context, _joinRoom);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        body: SingleChildScrollView(
            child: Container(
      padding: const EdgeInsets.all(20.0),
      child: Column(children: [
        Container(
          margin: const EdgeInsets.only(top: 20.0, bottom: 20.0),
          child: Row(
            children: [
              Expanded(
                  child: TextFormField(
                      decoration:
                          const InputDecoration(labelText: 'Search by code'),
                      onChanged: (value) {
                        setState(() {
                          roomFilter = value;
                        });
                      })),
              IconButton(
                onPressed: () {
                  _updateListRoomState();
                },
                icon: const Icon(
                  Icons.refresh,
                  color: Colors.pinkAccent,
                ),
                iconSize: 30.0,
              )
            ],
          ),
        ),
        SizedBox(
          height: 180,
          child: FutureBuilder<List<Room>>(
            future: rooms,
            builder: (context, snapshot) {
              if (snapshot.hasData) {
                final rooms = snapshot.data!
                    .where((room) => room.code
                        .toLowerCase()
                        .contains(roomFilter.toLowerCase()))
                    .toList();

                return GridView.builder(
                    gridDelegate:
                        const SliverGridDelegateWithFixedCrossAxisCount(
                            crossAxisCount: 8, mainAxisExtent: 130),
                    itemCount: rooms.length,
                    itemBuilder: (BuildContext ctx, index) {
                      return _RoomWidget(
                        id: rooms[index].id,
                        private: rooms[index].private,
                        code: rooms[index].code,
                        onTap: () {
                          _onTap(context, rooms[index].id);
                        },
                      );
                    });
              } else if (snapshot.hasError) {
                return const Text(
                    'Can not connect to the server. Check your internet connection!');
              }
              return const CircularProgressIndicator();
            },
          ),
        ),
        ElevatedButton(
            style: ButtonStyle(
                shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                    RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(20)))),
            onPressed: () async {
              _showAskingMoneyDialog(context, _newRoom);
            },
            child: const Text('New room')),
      ]),
    )));
  }
}

class _RoomWidget extends StatelessWidget {
  final int id;
  final String code;
  final bool private;
  final void Function() onTap;

  const _RoomWidget(
      {Key? key,
      this.id = 0,
      this.code = '',
      this.private = true,
      required this.onTap})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    AssetImage roomImage = private
        ? const AssetImage('assets/images/locked_room.png')
        : const AssetImage('assets/images/normal_room.png');
    return GestureDetector(
      child: Column(
        children: [
          Container(
            width: 50,
            height: 80,
            decoration: BoxDecoration(
              image: DecorationImage(image: roomImage),
            ),
          ),
          Text(
            code,
            style: const TextStyle(
                color: Colors.blue, fontWeight: FontWeight.bold),
          )
        ],
      ),
      onTap: onTap,
    );
  }
}

class _MySlider extends StatefulWidget {
  const _MySlider({this.maxValue = 0, required this.onChangedEndHandler});

  final double maxValue;
  final void Function(double) onChangedEndHandler;

  @override
  State<_MySlider> createState() => _MySliderState();
}

class _MySliderState extends State<_MySlider> {
  double _currentSliderValue = 0;

  @override
  Widget build(BuildContext context) {
    return Slider(
      value: _currentSliderValue,
      max: widget.maxValue,
      divisions: 100,
      label: _currentSliderValue.round().toString(),
      onChanged: (double value) {
        setState(() {
          _currentSliderValue = value;
        });
      },
      onChangeEnd: (double value) => {widget.onChangedEndHandler(value)},
    );
  }
}
