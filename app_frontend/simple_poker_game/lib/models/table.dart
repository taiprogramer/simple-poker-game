class _Card {
  final int id;
  final int number;
  final int suit;
  final String image;

  _Card(
      {required this.id,
      required this.number,
      required this.suit,
      required this.image});

  Map<String, dynamic> toJson() =>
      {'id': id, 'number': number, 'suit': suit, 'image': image};

  factory _Card.fromMap(Map data) {
    return _Card(
        id: data['id'],
        number: data['number'],
        suit: data['suit'],
        image: data['image']);
  }
}

class UserTurn {
  final int userID;

  UserTurn({this.userID = 0});

  Map<String, dynamic> toJson() => {'user_id': userID};

  factory UserTurn.fromMap(Map data) {
    return UserTurn(userID: data['user_id']);
  }
}

class PokerTable {
  final int id;
  final int round;
  final bool done;
  final int pot;
  final List<_Card> commonCards;
  final List<_Card> ownCards;
  final UserTurn currentTurn;

  PokerTable(
      {this.id = 0,
      this.round = 0,
      this.done = false,
      this.pot = 0,
      this.commonCards = const [],
      this.ownCards = const [],
      required this.currentTurn});

  Map<String, dynamic> toJson() => {
        'id': id,
        'round': round,
        'done': done,
        'pot': pot,
        'common_cards': commonCards,
        'own_cards': ownCards
      };
  factory PokerTable.fromMap(Map data) {
    return PokerTable(
        id: data['id'],
        round: data['round'],
        done: data['done'],
        pot: data['pot'],
        commonCards: _convertListDynamicToCards(data['common_cards']),
        ownCards: _convertListDynamicToCards(data['own_cards']),
        currentTurn: UserTurn.fromMap(data['current_turn']));
  }

  static List<_Card> _convertListDynamicToCards(List<dynamic> cardsData) {
    List<_Card> cards = List.empty(growable: true);
    for (final cardData in cardsData) {
      cards.add(_Card.fromMap(cardData));
    }
    return cards;
  }
}
