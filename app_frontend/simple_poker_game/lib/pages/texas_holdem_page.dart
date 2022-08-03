import 'package:flutter/material.dart';
import 'package:simple_poker_game/pages/room_texas_holdem_page.dart';
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

  @override
  void initState() {
    super.initState();
    rooms = RoomService.listRoom(0, 8);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text(widget.title),
        ),
        body: SingleChildScrollView(
            child: Container(
          padding: const EdgeInsets.all(20.0),
          child: Column(children: [
            const Center(
                child: Text(
              'Texas Hold\'em',
              style: TextStyle(
                  fontSize: 36.0,
                  fontWeight: FontWeight.bold,
                  color: Colors.pink),
            )),
            Container(
              margin: const EdgeInsets.only(top: 20.0, bottom: 20.0),
              child: Row(
                children: [
                  Expanded(
                      child: TextFormField(
                    decoration:
                        const InputDecoration(labelText: 'Search by code'),
                  )),
                  IconButton(
                    onPressed: () {},
                    icon: const Icon(
                      Icons.search,
                      color: Colors.pinkAccent,
                    ),
                    iconSize: 30.0,
                  )
                ],
              ),
            ),
            SizedBox(
              height: 300,
              child: FutureBuilder<List<Room>>(
                future: rooms,
                builder: (context, snapshot) {
                  if (snapshot.hasData) {
                    return GridView.builder(
                        gridDelegate:
                            const SliverGridDelegateWithFixedCrossAxisCount(
                                crossAxisCount: 4, mainAxisExtent: 100),
                        itemCount: snapshot.data!.length,
                        itemBuilder: (BuildContext ctx, index) {
                          return _RoomWidget(
                            private: snapshot.data![index].private,
                            code: snapshot.data![index].code,
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
                  int userID = AppLocalStorage.getItem("user_id");
                  RoomService.newRoom(userID: userID);
                },
                child: const Text('New room')),
          ]),
        )));
  }
}

class _RoomWidget extends StatelessWidget {
  final String code;
  final bool private;

  const _RoomWidget({Key? key, this.code = '', this.private = true})
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
      onTap: () {
        Navigator.pushReplacementNamed(context, RoomTexasHoldemPage.routeName);
      },
    );
  }
}
