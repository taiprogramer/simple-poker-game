import 'package:flutter/material.dart';

class RoomTexasHoldemPage extends StatelessWidget {
  static const String routeName = '/roomTexasHoldem';
  const RoomTexasHoldemPage({Key? key}) : super(key: key);

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
                  '75800',
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
                      _PlayerCircle(),
                      Container(
                        child: _PlayerCircle(),
                        margin: const EdgeInsets.only(bottom: 40.0),
                      ),
                      Container(
                        child: _PlayerCircle(),
                        margin: const EdgeInsets.only(bottom: 40.0),
                      ),
                      _PlayerCircle(),
                    ],
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _PlayerCircle(),
                      _PlayerCircle(),
                    ],
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _PlayerCircle(),
                      Container(
                        margin: const EdgeInsets.only(top: 40.0),
                        child: _PlayerCircle(),
                      ),
                      _PlayerCircle(),
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
    return Column(
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
              color: Colors.blue, borderRadius: BorderRadius.circular(100)),
          child: const Text(
            'G',
            style: TextStyle(fontWeight: FontWeight.bold),
          ),
        ),
        const Text(
          '\$ 999',
          style: TextStyle(color: Colors.red, backgroundColor: Colors.yellow),
        )
      ],
    );
  }
}
