MATCH (n)
DETACH DELETE n

// MATCH (p:Player)
// DELETE p;

MERGE (lo:Location {name: 'Lobby', description: 'A vast marble lined hall, this is the game lobby. You are safe here'}) 
RETURN lo;


MERGE (p:Player {username: 'user', name: 'bob', class: 'wizard', description: 'man'})-[:IN]->(l:Location {name: 'Lobby'})
RETURN p