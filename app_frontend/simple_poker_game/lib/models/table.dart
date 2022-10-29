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

class _Action {
  final int id;
  final String name;
  final int amount;
  _Action({this.id = 0, this.name = '', this.amount = 0});

  Map<String, dynamic> toJson() => {'id': id, 'name': name, 'amount': amount};

  factory _Action.fromMap(Map data) {
    return _Action(id: data['id'], name: data['name'], amount: data['amount']);
  }
}

class _BetHistory {
  final int userID;
  _Action action = _Action();

  _BetHistory({this.userID = 0, required this.action});

  factory _BetHistory.fromMap(Map data) {
    return _BetHistory(
        userID: data['user_id'], action: _Action.fromMap(data['action']));
  }
}

class _Result {
  final int userID;
  final List<_Card> cards;
  _Result({this.userID = 0, this.cards = const []});

  factory _Result.fromMap(Map data) {
    return _Result(
        userID: data['user_id'],
        cards: _convertListDynamicToCards(data['cards']));
  }

  static List<_Card> _convertListDynamicToCards(List<dynamic> cardsData) {
    List<_Card> cards = List.empty(growable: true);
    for (final cardData in cardsData) {
      cards.add(_Card.fromMap(cardData));
    }
    return cards;
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
  _BetHistory latestBet = _BetHistory(action: _Action());
  final List<_Result> results;

  PokerTable({
    this.id = 0,
    this.round = 0,
    this.done = false,
    this.pot = 0,
    this.commonCards = const [],
    this.ownCards = const [],
    this.results = const [],
    required this.currentTurn,
  });

  Map<String, dynamic> toJson() => {
        'id': id,
        'round': round,
        'done': done,
        'pot': pot,
        'common_cards': commonCards,
        'own_cards': ownCards,
        'results': results
      };
  factory PokerTable.fromMap(Map data) {
    final poker = PokerTable(
      id: data['id'],
      round: data['round'],
      done: data['done'],
      pot: data['pot'],
      commonCards: _convertListDynamicToCards(data['common_cards']),
      ownCards: _convertListDynamicToCards(data['own_cards']),
      currentTurn: UserTurn.fromMap(data['current_turn']),
      results: _convertListDynamicToResults(data['results']),
    );
    poker.latestBet = _BetHistory.fromMap(data['latest_bet']);
    return poker;
  }

  static List<_Card> _convertListDynamicToCards(List<dynamic> cardsData) {
    List<_Card> cards = List.empty(growable: true);
    for (final cardData in cardsData) {
      cards.add(_Card.fromMap(cardData));
    }
    return cards;
  }

  static List<_Result> _convertListDynamicToResults(List<dynamic> resultsData) {
    List<_Result> results = List.empty(growable: true);
    for (final result in resultsData) {
      results.add(_Result.fromMap(result));
    }
    return results;
  }
}
