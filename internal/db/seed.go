package db

import (
	"context"
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
}

var tags []string = []string{
	"productivity", "language-learning", "artificial-intelligence", "fitness", "travel",
	"public-speaking", "programming", "startup", "sleep", "creativity",
	"positive-thinking", "cooking", "investing", "social-media", "books",
	"relationships", "time-management", "healthcare", "reading", "motivation",
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
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user: ", err.Error())
			return
		}
	}

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

	log.Println("✅ Db seeded successfully.")
	return
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := range num {
		users[i] = &store.User{
			Username: names[rand.Intn(len(names))] + fmt.Sprintf("%d", i),
			Email:    names[rand.Intn(len(names))] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
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
