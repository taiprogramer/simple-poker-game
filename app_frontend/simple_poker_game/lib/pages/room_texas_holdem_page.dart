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
            height: 300,
            child: Container(
              color: Colors.red,
            ),
          ),
          Container(
            padding: const EdgeInsets.all(20.0),
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
            padding: const EdgeInsets.all(20.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                ElevatedButton(onPressed: () {}, child: const Text('Ready')),
              ],
            ),
          ),
          Container(
            padding: const EdgeInsets.all(20.0),
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
