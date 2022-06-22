package pg

type ChannelRepo struct {
	db *DB
}

func NewChannelRepo(db *DB) *ChannelRepo {
	return &ChannelRepo{db: db}
}

func (c *ChannelRepo) SayHello() string {
	return "hello world"
}
