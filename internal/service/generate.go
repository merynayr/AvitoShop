package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i ShopService -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i AuthService -o ./mocks/ -s "_minimock.go"
