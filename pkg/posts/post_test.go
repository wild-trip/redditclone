package posts

import (
	"fmt"
	"reddit/pkg/user"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	testPost1 = &Post{
		Author: &user.User{
			ID:       1,
			Username: "rvasily",
		},
		Category: "music",
		CommentsID: []bson.ObjectId{
			"^\xba\xf9\xf2<\x04\xc1|V\xf5\x12E",
			"^\xbb\xfes<\x04\xc1+\x9b\xf8]\x1b",
		},
		Created:          "2020-05-12T22:33:02+03:00",
		ID:               "^\xba\xf9\xee<\x04\xc1|V\xf5\x12D",
		Score:            1,
		Text:             "Something 1",
		Link:             "",
		Title:            "Lorem",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            10,
		Votes: []Vote{
			{
				UserID: 1,
				Rating: 1,
			},
		},
	}
	testPost2 = &Post{
		Author: &user.User{
			ID:       1,
			Username: "rvasily",
		},
		Category: "funny",
		CommentsID: []bson.ObjectId{
			"^\xba\xf9\xf2<\x04\xc1|V\xf5\x12E",
			"^\xbb\xfes<\x04\xc1+\x9b\xf8]\x1b",
		},
		Created:          "2020-05-12T22:33:02+03:00",
		ID:               "^\xba\xf9\xee<\x04\xc1|V\xf5\x12D",
		Score:            1,
		Text:             "Something 1",
		Link:             "",
		Title:            "lol",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            5,
		Votes: []Vote{
			{
				UserID: 1,
				Rating: 1,
			},
		},
	}
	testPost3 = &Post{
		Author: &user.User{
			ID:       1,
			Username: "rvasily",
		},
		Category:         "funny",
		CommentsID:       []bson.ObjectId{},
		Created:          "2020-05-12T22:33:02+03:00",
		ID:               "^\xba\xf9\xee<\x04\xc1|V\xf5\x12D",
		Score:            1,
		Text:             "Something 1",
		Link:             "",
		Title:            "lol",
		Type:             "text",
		UpvotePercentage: 100,
		Views:            5,
		Votes: []Vote{
			{
				UserID: 1,
				Rating: 1,
			},
		},
	}
)

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := []*Post{testPost1, testPost2}

	mockDB.EXPECT().Find(bson.M{}).Return(mockDBFind)
	var responsePosts []*Post
	mockDBFind.EXPECT().All(gomock.Any()).SetArg(0, expectPosts)

	responsePosts, err := testRepo.GetAll()
	assert.Equal(t, expectPosts, responsePosts)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//BD error
	mockDB.EXPECT().Find(bson.M{}).Return(mockDBFind)
	mockDBFind.EXPECT().All(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responseErrPosts, err := testRepo.GetAll()

	assert.Empty(t, responseErrPosts)
	assert.EqualError(t, err, "DB err: Internal error")
}

func TestGetCategory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := []*Post{testPost2}
	category := testPost2.Category
	mockDB.EXPECT().Find(bson.M{"category": category}).Return(mockDBFind)
	var responsePosts []*Post
	mockDBFind.EXPECT().All(gomock.Any()).SetArg(0, expectPosts)

	responsePosts, err := testRepo.GetCategory(category)
	assert.Equal(t, expectPosts, responsePosts)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//BD error
	mockDB.EXPECT().Find(bson.M{"category": category}).Return(mockDBFind)
	mockDBFind.EXPECT().All(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responseErrPosts, err := testRepo.GetCategory(category)

	assert.Empty(t, responseErrPosts)
	assert.EqualError(t, err, "Internal error")
}

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost2
	postID := testPost2.ID
	mockDB.EXPECT().Find(bson.M{"_id": postID}).Return(mockDBFind)
	var responsePosts *Post
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	responsePosts, err := testRepo.GetByID(postID)
	assert.Equal(t, expectPosts, responsePosts)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//BD error
	mockDB.EXPECT().Find(bson.M{"_id": postID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responseErrPosts, err := testRepo.GetByID(postID)

	assert.Empty(t, responseErrPosts)
	assert.EqualError(t, err, "Internal error")
}

func TestGetByUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := []*Post{testPost1, testPost2}
	userLogin := testPost2.Author.Username
	mockDB.EXPECT().Find(bson.M{"author.username": userLogin}).Return(mockDBFind)
	var responsePosts []*Post
	mockDBFind.EXPECT().All(gomock.Any()).SetArg(0, expectPosts)

	responsePosts, err := testRepo.GetByUserLogin(userLogin)
	assert.Equal(t, expectPosts, responsePosts)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//BD error
	mockDB.EXPECT().Find(bson.M{"author.username": userLogin}).Return(mockDBFind)
	mockDBFind.EXPECT().All(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responseErrPosts, err := testRepo.GetByUserLogin(userLogin)

	assert.Empty(t, responseErrPosts)
	assert.EqualError(t, err, "Internal error")
}

func TestAddPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	//mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	user := testPost3.Author
	category := testPost3.Category
	title := testPost3.Title
	typePost := testPost3.Type
	text := testPost3.Text
	link := testPost3.Link

	expectPosts := testPost3
	//Add text post
	mockDB.EXPECT().Insert(gomock.Any()).Return(nil)
	var responsePost *Post
	responsePost, err := testRepo.Add(user, category, title, typePost, text, link)
	assert.Equal(t, expectPosts.Author, responsePost.Author)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Add link post
	typePost = "link"
	text = ""
	link = "www.google.com"
	mockDB.EXPECT().Insert(gomock.Any()).Return(nil)
	responsePost, err = testRepo.Add(user, category, title, typePost, text, link)
	assert.Equal(t, expectPosts.Author, responsePost.Author)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//BD error
	mockDB.EXPECT().Insert(gomock.Any()).Return(fmt.Errorf("Internal error"))
	responsePost, err = testRepo.Add(user, category, title, typePost, text, link)

	assert.Empty(t, responsePost)
	assert.EqualError(t, err, "Internal error")
}

func TestUpdatePostByNewComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost3
	commentID := bson.NewObjectId()

	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)
	var responsePost *Post

	responsePost, err := testRepo.AddComment(expectPosts.ID, commentID)
	assert.Equal(t, expectPosts.Author, responsePost.Author)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err get by id
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responsePostErr, err := testRepo.AddComment(expectPosts.ID, commentID)
	assert.Empty(t, responsePostErr)
	assert.EqualError(t, err, "Error getting post from BD: Internal error")

	//Err AddComment
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(fmt.Errorf("Internal server error"))

	responsePostErr2, err := testRepo.AddComment(expectPosts.ID, commentID)
	assert.Empty(t, responsePostErr2)
	assert.EqualError(t, err, "Error update BD: Internal server error")
}

func TestUpdatePostByUpViews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost3

	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)

	err := testRepo.UpViews(expectPosts.ID)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err get by id
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).Return(fmt.Errorf("Internal error"))

	err = testRepo.UpViews(expectPosts.ID)
	assert.EqualError(t, err, "Error getting post from BD: Internal error")

	//Err AddComment
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(fmt.Errorf("Internal server error"))

	err = testRepo.UpViews(expectPosts.ID)
	assert.EqualError(t, err, "Error update BD: Internal server error")

}

func TestUpdatePostByDeleteComment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost1
	commentID := expectPosts.CommentsID[0]

	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)
	var responsePost *Post

	responsePost, err := testRepo.DeleteComment(expectPosts.ID, commentID)
	assert.Equal(t, expectPosts.Author, responsePost.Author)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err get by id
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responsePostErr, err := testRepo.DeleteComment(expectPosts.ID, commentID)
	assert.Empty(t, responsePostErr)
	assert.EqualError(t, err, "Error getting post from BD: Internal error")

	//Err AddComment
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(fmt.Errorf("Internal server error"))

	responsePostErr2, err := testRepo.DeleteComment(expectPosts.ID, commentID)
	assert.Empty(t, responsePostErr2)
	assert.EqualError(t, err, "Error update BD: Internal server error")
}

func TestUpdatePostByUpvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost1
	user := &user.User{
		ID:       8,
		Username: "igor",
	}
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)

	responsePost, err := testRepo.Upvote(user, expectPosts.ID)

	assert.Equal(t, expectPosts.Category, responsePost.Category)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Add downwote post
	expectPosts.Votes[0].Rating = -1
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)

	responsePost, err = testRepo.Upvote(expectPosts.Author, expectPosts.ID)

	assert.Equal(t, expectPosts.Category, responsePost.Category)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err get by id
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responsePostErr, err := testRepo.Upvote(expectPosts.Author, expectPosts.ID)
	assert.Empty(t, responsePostErr)
	assert.EqualError(t, err, "DB err: Internal error")

	//Err AddComment
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(fmt.Errorf("Internal server error"))

	responsePostErr2, err := testRepo.Upvote(expectPosts.Author, expectPosts.ID)

	assert.Empty(t, responsePostErr2)
	assert.EqualError(t, err, "Error update BD: Internal server error")
}

func TestUpdatePostByDownvote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	mockDBFind := NewMockFindInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost1
	user := &user.User{
		ID:       9,
		Username: "lera",
	}
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)

	responsePost, err := testRepo.Downvote(user, expectPosts.ID)

	assert.Equal(t, expectPosts.Category, responsePost.Category)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Add downwote post
	expectPosts.Votes[0].Rating = 1
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(nil)

	responsePost, err = testRepo.Downvote(expectPosts.Author, expectPosts.ID)

	assert.Equal(t, expectPosts.Category, responsePost.Category)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err get by id
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).Return(fmt.Errorf("Internal error"))

	responsePostErr, err := testRepo.Downvote(expectPosts.Author, expectPosts.ID)
	assert.Empty(t, responsePostErr)
	assert.EqualError(t, err, "DB err: Internal error")

	//Err AddComment
	mockDB.EXPECT().Find(bson.M{"_id": expectPosts.ID}).Return(mockDBFind)
	mockDBFind.EXPECT().One(gomock.Any()).SetArg(0, expectPosts)

	mockDB.EXPECT().Update(bson.M{"_id": expectPosts.ID}, expectPosts).Return(fmt.Errorf("Internal server error"))

	responsePostErr2, err := testRepo.Downvote(expectPosts.Author, expectPosts.ID)

	assert.Empty(t, responsePostErr2)
	assert.EqualError(t, err, "Error update BD: Internal server error")
}

func TestDeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := NewMockPostRepositoryDBInterface(ctrl)
	testRepo := NewRepo(mockDB)

	expectPosts := testPost1

	mockDB.EXPECT().Remove(bson.M{"_id": expectPosts.ID}).Return(nil)

	isDelete, err := testRepo.Delete(expectPosts.ID)
	assert.True(t, isDelete)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err not found
	mockDB.EXPECT().Remove(bson.M{"_id": expectPosts.ID}).Return(mgo.ErrNotFound)

	isDelete, err = testRepo.Delete(expectPosts.ID)
	assert.False(t, isDelete)
	assert.Empty(t, err, fmt.Sprintf("Unexpected error: %v", err))

	//Err  BD error
	mockDB.EXPECT().Remove(bson.M{"_id": expectPosts.ID}).Return(fmt.Errorf("Internal error"))

	isDelete, err = testRepo.Delete(expectPosts.ID)
	assert.False(t, isDelete)
	assert.EqualError(t, err, "Internal error")
}
