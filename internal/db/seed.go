package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/mustaphalimar/go-social/internal/store"
)

var names []string = []string{
	"john", "emma", "michael", "sophia", "james", "olivia", "william", "isabella", "alexander", "mia",
	"daniel", "charlotte", "henry", "amelia", "joseph", "emily", "david", "harper", "benjamin", "evelyn",
	"lucas", "abigail", "matthew", "ella", "samuel", "avery", "christopher", "scarlett", "jack", "grace",
	"andrew", "chloe", "ethan", "lily", "ryan", "nora", "joshua", "hannah", "nathan", "zoe",
	"caleb", "madison", "sebastian", "layla", "elijah", "victoria", "gabriel", "penelope", "dylan", "aria",
	"leo", "samantha", "julian", "hazel", "isaac", "ellie", "mason", "aurora", "logan", "stella",
	"jackson", "paisley", "aiden", "savannah", "carter", "brooklyn", "wyatt", "bella", "owen", "claire",
	"peter", "lucy", "noah", "zoey", "mila", "caroline", "adam", "ruby", "max", "ivy",
	"juliet", "simon", "elise", "arthur", "faith", "riley", "sienna", "miles", "luna", "oliver",
}

var titles []string = []string{
	"10 Tips for a Productive Morning",
	"How to Learn a New Language Fast",
	"The Future of Artificial Intelligence",
	"Why Exercise is Key to a Healthy Life",
	"Top 5 Travel Destinations This Year",
	"Mastering the Art of Public Speaking",
	"Best Programming Languages to Learn in 2025",
	"How to Build a Successful Startup",
	"The Science Behind Good Sleep",
	"Simple Habits to Boost Your Creativity",
	"The Power of Positive Thinking",
	"How to Cook the Perfect Steak",
	"Investing 101: Where to Start",
	"The Psychology of Social Media",
	"Top 10 Books That Will Change Your Life",
	"Secrets to a Happy Relationship",
	"How to Master Time Management",
	"The Role of AI in Modern Healthcare",
	"Why Reading is Essential for Growth",
	"How to Stay Motivated Every Day",
	"How to Meditate for Beginners",
	"Budget-Friendly Meal Prep Ideas",
	"Essential Skills for Remote Work",
	"How to Start a Podcast",
	"Gardening Tips for Small Spaces",
	"How to Improve Your Memory",
	"Traveling Safely in 2025",
	"Building Emotional Intelligence",
	"How to Write Engaging Content",
	"Tips for Sustainable Living",
}

var contents []string = []string{
	"Start your day right with these proven productivity hacks.",
	"Master a new language quickly with these expert tips.",
	"Exploring how AI is shaping the future of technology and business.",
	"Discover why regular exercise is essential for long-term health.",
	"The top must-visit places to add to your travel bucket list.",
	"Learn how to speak confidently and captivate any audience.",
	"An in-depth look at the best programming languages in 2025.",
	"Step-by-step guide to launching and scaling a successful startup.",
	"Unveiling the science behind quality sleep and its benefits.",
	"Boost your creativity with these simple and effective habits.",
	"How adopting a positive mindset can transform your life.",
	"A foolproof method for cooking the juiciest steak every time.",
	"A beginner-friendly guide to smart investing strategies.",
	"How social media impacts mental health and relationships.",
	"Books that will inspire, educate, and change your perspective.",
	"Discover the key elements that make relationships last.",
	"Simple strategies to manage time and increase daily efficiency.",
	"How AI is revolutionizing healthcare and medical treatments.",
	"Unlock the benefits of daily reading for personal growth.",
	"Practical techniques to stay motivated and achieve your goals.",
	"Learn the basics of meditation and how to get started today.",
	"Delicious and affordable meal prep ideas for busy people.",
	"Key skills you need to thrive while working remotely.",
	"Everything you need to know to launch your own podcast.",
	"Grow your own food with these small-space gardening tips.",
	"Proven methods to boost your memory and recall.",
	"Stay safe while traveling with these up-to-date tips.",
	"Develop your emotional intelligence for better relationships.",
	"Write content that captures attention and drives engagement.",
	"Easy ways to live a more sustainable and eco-friendly life.",
}

var tags []string = []string{
	"productivity", "language-learning", "artificial-intelligence", "fitness", "travel",
	"public-speaking", "programming", "startup", "sleep", "creativity",
	"positive-thinking", "cooking", "investing", "social-media", "books",
	"relationships", "time-management", "healthcare", "reading", "motivation",
	"meditation", "meal-prep", "remote-work", "podcasting", "gardening",
	"memory", "safety", "emotional-intelligence", "writing", "sustainability",
}

var comments []string = []string{
	"Great tips, thanks!",
	"Love this! Very helpful.",
	"AI is the future!",
	"Need to exercise more.",
	"Adding these places to my list!",
	"Public speaking tips are gold.",
	"Go is definitely on my list.",
	"Starting my business soon, great advice!",
	"Sleep science is fascinating.",
	"Love these creativity habits.",
	"Positive thinking really works.",
	"Perfect steak tips!",
	"Investing feels overwhelming.",
	"Social media is powerful.",
	"Great book recommendations!",
	"Communication is key in relationships.",
	"Time management is a struggle.",
	"AI is a game changer for healthcare.",
	"Reading daily makes a difference.",
	"Staying motivated is hard.",
	"Just what I needed to read today.",
	"These meal prep ideas are awesome!",
	"Remote work can be tough, thanks for the advice.",
	"Starting a podcast has been on my mind.",
	"Small-space gardening is so rewarding.",
	"Memory tips are super useful.",
	"Travel safety is more important than ever.",
	"Emotional intelligence is underrated.",
	"Writing engaging content is a challenge.",
	"Sustainability tips are always welcome.",
}

// ... your imports and existing code ...

// Add more data to your arrays as shown above

func EraseAll(store store.Storage) {
	ctx := context.Background()
	if err := store.Comments.DeleteAll(ctx); err != nil {
		log.Println("Error deleting comments:", err)
	}
	if err := store.Posts.DeleteAll(ctx); err != nil {
		log.Println("Error deleting posts:", err)
	}
	if err := store.Users.DeleteAll(ctx); err != nil {
		log.Println("Error deleting users:", err)
	}
	log.Println("🗑️  All data erased from the database.")
}

func Seed(store store.Storage, db *sql.DB) {
	EraseAll(store)
	ctx := context.Background()

	users := generateUsers(100)
	tx, _ := db.BeginTx(ctx, nil)
	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user: ", err.Error())
			return
		}
	}
	tx.Commit()

	posts := generatePosts(50, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating Post: ", err.Error())
			return
		}
	}

	comments := generateComments(100, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating Comment: ", err.Error())
			return
		}
	}
	followers := generateFollowers(150, users) // e.g., 150 follow relationships
	for _, f := range followers {
		if err := store.Followers.Follow(ctx, f.FollowerId, f.UserId); err != nil {
			log.Println("Error creating Follower: ", err.Error())
			return
		}
	}

	log.Println("✅ Db seeded successfully. 💯")
	return
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		users[i] = &store.User{
			Username: names[rand.Intn(len(names))] + fmt.Sprintf("%d", i),
			Email:    names[rand.Intn(len(names))] + fmt.Sprintf("%d", i) + "@example.com",
			RoleID:   1,
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := range num {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: titles[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}
	return posts

}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	randComments := make([]*store.Comment, num)

	for i := range num {

		randComments[i] = &store.Comment{
			UserID:  users[rand.Intn(len(users))].ID,
			PostID:  posts[rand.Intn(len(posts))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return randComments
}

func generateFollowers(num int, users []*store.User) []*store.Follower {
	followers := make([]*store.Follower, 0, num)
	userCount := len(users)
	seen := make(map[string]struct{})

	for len(followers) < num {
		followerIdx := rand.Intn(userCount)
		followedIdx := rand.Intn(userCount)
		if followerIdx == followedIdx {
			continue // no self-follow
		}
		key := fmt.Sprintf("%d-%d", users[followerIdx].ID, users[followedIdx].ID)
		if _, exists := seen[key]; exists {
			continue // no duplicate follows
		}
		seen[key] = struct{}{}
		followers = append(followers, &store.Follower{
			FollowerId: users[followerIdx].ID,
			UserId:     users[followedIdx].ID,
		})
	}
	return followers
}
