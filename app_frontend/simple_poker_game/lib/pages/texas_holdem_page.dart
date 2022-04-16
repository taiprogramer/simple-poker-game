import 'package:flutter/material.dart';

class TexasHoldemPage extends StatelessWidget {
  static const String routeName = '/texasHoldem';

  const TexasHoldemPage({Key? key, this.title = 'Simple Poker Game'})
      : super(key: key);

  final String title;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text(title),
        ),
        body: SingleChildScrollView(
            child: Container(
          padding: const EdgeInsets.all(20.0),
          child: Column(children: [
            const Center(
                child: Text(
              'Texas Hold\'em',
              style: TextStyle(fontSize: 36.0, fontWeight: FontWeight.bold),
            )),
            Container(
              margin: const EdgeInsets.only(top: 20.0, bottom: 20.0),
              child: Row(
                children: [
                  Expanded(
                      child: TextFormField(
                    decoration: const InputDecoration(
                        border: OutlineInputBorder(),
                        labelText: 'Search by code'),
                  )),
                  IconButton(
                    onPressed: () {},
                    icon: const Icon(Icons.search),
                    iconSize: 30.0,
                  )
                ],
              ),
            ),
            SizedBox(
                height: 300,
                child: SingleChildScrollView(
                    child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: const [
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(
                          private: false,
                        )
                      ],
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: const [
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(
                          private: false,
                        )
                      ],
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: const [
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(
                          private: false,
                        )
                      ],
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceAround,
                      children: const [
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(),
                        _RoomWidget(
                          private: false,
                        )
                      ],
                    ),
                  ],
                ))),
            ElevatedButton(onPressed: () {}, child: const Text('New room'))
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
          const Text(
            '#FAK0',
            style: TextStyle(color: Colors.blue, fontWeight: FontWeight.bold),
          )
        ],
      ),
      onTap: () {},
    );
  }
}
