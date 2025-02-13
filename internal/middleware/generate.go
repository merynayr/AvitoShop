package middleware

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i UserMiddleware -o ./mocks/ -s "_minimock.go"
