# systemd

First clone the `wof-staticd` unit file:

```
cp wof-staticd.service.example wof-staticd.service
```

Adjust the file to taste. It assumes as few things that you may need or want to tweak:

* That you've built and copied go-whosonfirst-static/bin/wof-staticd to `/usr/local/bin/wof-staticd`
* That you want to run `wof-staticd` as user `www-data` (you will need to change this if you're running CentOS for example)
* That you will replace the `-stuff -stuff -stuff` flags in the `ExecStart` with meaningful config flags

Move the file in to place:

```
mv wof-staticd.service /lib/systemd/system/wof-staticd.service
```

Now tell `systemd` about it:

```
systemctl enable wof-staticd.service
sudo systemctl start wof-staticd
```

## See also

* https://fabianlee.org/2017/05/21/golang-running-a-go-binary-as-a-systemd-service-on-ubuntu-16-04/