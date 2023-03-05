module github.com/bjarne-dietrich/pingo

go 1.19

require internal/icmp v1.0.0
replace internal/icmp => ./internal/icmp

require internal/utils v1.0.0
replace internal/utils => ./internal/utils
