package neo4j

// Now this is for Neo4J.
var queries = map[string]string{
	// Add Version:
	"releaseVersion": `
		CREATE (v:Version { version: {version}, forced: {forced} })
		RETURN v;
	`,
	// Check Version:
	"getAllVersion": `
		MATCH (ver:Version)
		WHERE ver.forced = TRUE
		RETURN ver ORDER BY ver.version DESC;
	`,
	"getAllForcedVersion": `
		MATCH (v:Version)
		RETURN v ORDER BY v.version DESC;
	`,
	// SMS Verification:
	"mergeVerificationRequest": `
		MERGE (vr:Verification { phoneNumber: {phoneNumber} })
		SET 
			vr.code        = {code},
			vr.token       = {token},
			vr.verified    = FALSE, 
			vr.ttl         = timestamp();
	`,
	"verifyRequest": `
		MATCH (vr:Verification)
		WHERE vr.phoneNumber = {phoneNumber} AND vr.code = {code}
		SET vr.verified = TRUE
		RETURN vr.token;
	`,
	"isVerified": `
		MATCH (vr:Verification)
		WHERE vr.token = {token}
		RETURN vr.verified, vr.phoneNumber;
	`,
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
		RETURN device.deviceToken, user IS NOT NULL AS signedIn;
	`,
	// TODO ******
	"whoAmI": `
		MATCH (d:Device)-[:SIGNED]-(u:User)
		WHERE d.deviceToken = {deviceToken}
		RETURN u.username;
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
	// TODO
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
			p.follows    = 0

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
				u.fullName   = "",
				u.bio        = "",
				u.location   = "",
				u.followedBy = 0,
				u.follows    = 0
		SET p += {change};
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
		FOREACH (x IN seconds | CREATE (lu)-[:Next]->(x));
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