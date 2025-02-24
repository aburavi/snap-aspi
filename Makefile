.PHONY: run stop

run:
	./bin/snap_auth &
	./bin/snap_backend &
	./bin/snap_signature &
	./bin/snap_ratelimiter &
	./bin/snap_history &
	./bin/snap_inquiry &
	./bin/snap_transfer &
	./bin/snap_historyv2 &
	./bin/snap_inquiryv2 &
	./bin/snap_transferv2 &
	./bin/gateway &

stop:
	killall snap_auth
	killall snap_backend
	killall snap_signature
	killall snap_ratelimiter
	killall snap_history
	killall snap_inquiry
	killall snap_transfer
	killall snap_historyv2
	killall snap_inquiryv2
	killall snap_transferv2
	killall gateway
	
