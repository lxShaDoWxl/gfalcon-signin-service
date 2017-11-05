# SignIn Service

# Run
```bash
$ dep ensure
$ cd vendor/github.com/m0cchi/gfalcon/metatools/sql
$ mysql -u $USER -p < create_database.sql
$ mysql -u $USER -p -D gfalcon < create_database.sql
$ mysql -u $USER -p -D gfalcon < create_table.sql
$ mysql -u $USER -p -D gfalcon < create_defaultdata.sql
$ cd ../../../../../../
$ bower install
$ go run test/init_data.go --dbhost 'user:password@unix(/tmp/mysql.sock)/gfalcon?parseTime=true'
$ go run app/server.go --dbhost 'user:password@unix(/tmp/mysql.sock)/gfalcon?parseTime=true'
```

# feature
This product is SingIn Service.
If authentication succeeds, user got `gfalcon.session` in cookie.
SP verifies whether it has been authenticated using `gfalcon.session`.

# TODO
- SSO View
- show service list

# View
![](https://i.gyazo.com/46b429a45ee882638cf92c90e0da251e.gif)
