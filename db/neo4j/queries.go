package neo4j

// Now this is for Neo4J.
var queries = map[string]string{
	// Auth:
	"touchDevice": `
		MERGE (device:Device { uid: {uid} })
		ON CREATE SET
			device.deviceToken = {maybeDeviceToken}, 
			device.uid = {uid}, 
			device.name = {name},
			device.platform = {platform},
			device.capacity = {capacity},
			device.os_type = {os_type},
			device.os_version = {os_version}
		WITH device
		OPTIONAL MATCH (device)-[:SIGNED]->(user:User)
		RETURN device.deviceToken, user IS NOT NULL AS signedIn, user.username;
	`,
	"ensureDeviceToken": `
		MATCH (d:Device)
		WHERE d.deviceToken = {deviceToken}
		RETURN d.deviceToken;
	`,
	"whoAmI": `
		MATCH (d:Device)-[:SIGNED]-(u:User)
		WHERE d.deviceToken = {deviceToken}
		RETURN u.username;
	`,
	"isUniquePhoneNumber": `
		MATCH (u:User)
		WHERE u.phoneNumber = {phoneNumber}
		RETURN u;
	`,
	"isUniqueUsername": `
		MATCH (u:User)
		WHERE u.username = {username}
		RETURN u;
	`,
	"signDeviceIn": `
		MATCH (d:Device) WHERE d.deviceToken = {deviceToken} 
		WITH d
		MATCH (u:User) WHERE u.username = {username} AND u.password = {password}
		MERGE (d)-[:SIGNED]->(u)
		RETURN d, u;
	`,
	"signDeviceOut": `
		MATCH (d:Device)
		WITH d
		OPTIONAL MATCH (d)-[s:SIGNED]-(:User)
		WHERE d.deviceToken = {deviceToken}
		DELETE s
		RETURN d;
	`,
	"signUp": `
		MATCH (d:Device) WHERE d.deviceToken = {deviceToken}
		CREATE (d)-[:SIGNED]->(u:User)-[:BIND]->(p:Profile)
		SET
			u.username = {username},
			u.password = {password},
			u.phoneNumber = {phoneNumber},

			p.fullName   = "",
			p.bio        = "",
			p.location   = "",
			p.followedBy = 0,
			p.follows    = 0,
			p.avatar     = "",
			p.wallpaper  = ""

		RETURN d;
		`,
	// Profile:
	"getProfile": `
		MATCH (u:User)-[:BIND]-(p:Profile)
		WHERE u.username = {username}
		RETURN p;
	`,
	"updateProfile": `
		MATCH (u:User) WHERE u.username = {username}
		MERGE (u)-[:BIND]-(p:Profile)
			ON CREATE SET 
				p.fullName   = "",
				p.bio        = "",
				p.location   = "",
				p.followedBy = 0,
				p.follows    = 0,
				p.avatar     = "",
				p.wallpaper  = ""
		SET p += {change}
		RETURN u;
	`,
	"isFollowedBy": `
		OPTIONAL MATCH (u1:User)-[:BIND]-(p1:Profile)-[f:FOLLOW]->(p2:Profile)-[:BIND]-(u2:User)
		WHERE u1.username = {username1} AND u2.username = {username2}
		RETURN f IS NOT NULL;
	`,
	"follow": `
		MATCH (u1:User)-[:BIND]-(p1:Profile), (u2:User)-[:BIND]-(p2:Profile)
		WHERE u1.username = {username1} AND u2.username = {username2}
		MERGE (p1)-[f:FOLLOW]->(p2)
			ON CREATE SET
				f.created_at  = timestamp(),
				p1.follows    = p1.follows + 1,
				p2.followedBy = p1.followedBy + 1

		RETURN TRUE;
	`,
	"unfollow": `
		MATCH (u1:User)-[:BIND]->(p1:Profile)-[f:FOLLOW]->(p2:Profile)<-[:BIND]-(u2:User)
		WHERE u1.username = {username1} AND u2.username = {username2}
		DELETE f
		WITH p1, p2
		SET 
			p1.follows    = p1.follows - 1, 
			p2.followedBy = p1.followedBy - 1
		RETURN TRUE;
	`,
	"post": `
		MATCH (u:User)-[:BIND]-(p:Profile) WHERE u.username = {username}
		OPTIONAL MATCH (p)-[r:Post]-(secondlatestupdate)
		DELETE r
		CREATE (p)-[:Post]->(lu:Post)
		WITH lu, secondlatestupdate
		SET
			lu.artID = {artID},
			lu.title = {title},
			lu.desc  = {desc},
			lu.likes_count = 0,
			lu.comments_count = 0,
			lu.tags = {tags},
			lu.date = timestamp(),
			lu.displaySource = {displaySource}
		WITH lu, collect(secondlatestupdate) AS seconds
		FOREACH (x IN seconds | CREATE (lu)-[:Next]->(x))
		RETURN lu;
	`,
	"like": `
		MATCH (u:User)-[:BIND]-(profile:Profile) WHERE u.username = {username}
		WITH profile
		MATCH (post:Post) WHERE post.artID = {artID}
		MERGE (profile)-[r1:OWN]->(like:Like)-[r2:THAT]->(post)
			ON CREATE SET
				post.likes_count = post.likes_count + 1,
				r1.created_at = timestamp(),
				r2.created_at = timestamp(),
				like.created_at = timestamp(),
				like.flag = 1 // for if like statement
		
		WITH post, profile, like

		MATCH (like { flag: 1 })
		OPTIONAL MATCH (post)-[r:LIKED_BY]-(secondLatestUpdate)
		DELETE r
		CREATE (post)-[:LIKED_BY]->(like)
		WITH like, collect(secondLatestUpdate) AS seconds
		FOREACH (x IN seconds | CREATE (like)-[:NEXT]->(x))
		REMOVE like.flag
		RETURN like
	`,
	"dislike": `
		MATCH (u:User)-[:BIND]-(prof:Profile) WHERE u.username = {username}
		WITH prof
		MATCH (post:Post) WHERE post.artID = {artID}
		RETURN TRUE
		
		UNION

		MATCH (prof)-[:OWN]-(like:Like)-[:THAT]-(post)
		SET post.likes_count = post.likes_count - 1

		WITH like

		OPTIONAL MATCH (like)-[:NEXT]->(next)
		OPTIONAL MATCH (first)-[:LIKE_BY]->(like)
		OPTIONAL MATCH (prev)-[:NEXT]->(like)
        
        WITH like, collect(next) AS next, collect(first) AS first, collect(prev) AS prev

		FOREACH (x IN next | 
			FOREACH (y IN first | CREATE (y)-[:LIKED_BY]->(x))
		)
		FOREACH (x IN next | 
			FOREACH (y IN prev | CREATE (y)-[:NEXT]->(x))
		)
		DETACH DELETE like

		RETURN TRUE;
	`,
	"getLikes": `
		MATCH (post:Post) WHERE post.artID = {artID}
		MATCH (post)-[:LIKED_BY]-(start:Like)-[:NEXT*0..]-(like)-[:OWN]-(:Profile)-[:BIND]-(user:User) WHERE like.created_at < {cursur}
		RETURN user.username, like.created_at ORDER BY like.created_at DESC LIMIT {count};
	`,
	"getPosts": `
		MATCH (u:User)-[:BIND]-(p:Profile) WHERE u.username = {username}
		WITH p
		MATCH (p)-[:Post]-(start)-[:Next*0..]-(post) WHERE post.date < {cursur}
		RETURN post ORDER BY p.date DESC LIMIT {count};
	`,
	// WARNING: max depth is 20
	"getTimeline": `
		MATCH (u:User)-[:BIND]-(p:Profile) WHERE u.username = {username}
		WITH p
		MATCH (p)-[:FOLLOW*0..1]->(f:Profile)
		WITH f
		MATCH (f)-[:Post]-(start)-[:Next*0..20]-(post) WHERE post.date < {cursur}
		RETURN post ORDER BY post.date DESC LIMIT {count};
	`,
}

func (db *neo4jDB) GetQuery(key string) interface{} {
	if value, ok := queries[key]; ok {
		return value
	}
	panic("query " + key + " does not exist")
}
