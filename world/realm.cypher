MATCH (n)
DETACH DELETE n;

// ------------------

CREATE (lobby:Location {name: 'Lobby', description: 'A vast marble lined hall, this is the game lobby. You are safe here'});

// ------------------

MATCH (lobby:Location {name: 'Lobby'})
MERGE (p:Player {username: 'bob@bob.com', name: 'Jimmy Guy', class: 'Wizard', description: 'Hairy'})-[:IN]->(lobby)
RETURN p;

// ------------------

MATCH (lobby:Location {name: 'Lobby'})
MERGE (p:Player {username: 'fred@bob.com'})-[:IN]->(lobby)
ON CREATE
  SET p.description = 'Green and hairy', p.class = 'Rogue', p.name = 'Freddy'
ON MATCH
  SET p.description = 'Green and hairy', p.class = 'Rogue', p.name = 'Freddy'
RETURN p;

MATCH (lobby:Location {name: 'Lobby'})
MERGE (p:Player {username: 'fred@bob.com'})-[:IN]->(lobby)
ON CREATE
  SET p.description = 'Green and hairy', p.class = 'WANK', p.name = 'Freddy'
ON MATCH
  SET p.description = 'Green and hairy', p.class = 'WANK', p.name = 'Freddy'
RETURN p;

// MERGE (p:Player {username: 'fred@bob.com'})
// ON MATCH
//   SET p.description = 'Green and hairy', p.class = 'Rogue', p.name = 'Freddy';

// ------------------

MATCH (l:Location {name: 'Lobby'})<--(playersInLobby:Player)
RETURN playersInLobby