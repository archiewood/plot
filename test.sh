# bar chart
go run . -t "Bar: Big Numbers" test/big_numbers.csv
go run . -t "Bar: Invalid Second Column" test/invalid_second_column.csv
go run . -t "Bar: Long Labels" test/long_labels.csv
go run . -t "Bar: MultiTimeline" test/multiTimeline.csv
go run . -t "Bar: Single Column Header" test/single_col_header.csv
go run . -t "Bar: Single Column" test/single_col.csv
go run . -t "Bar: Two Columns" test/two_col.csv
go run . -t "Bar: Two Columns Header" test/two_col_header.csv
# column chart
go run . -t "Col: Big Numbers" -c column test/big_numbers.csv
go run . -t "Col: Invalid Second Column" -c column test/invalid_second_column.csv
go run . -t "Col: Long Labels" -c column test/long_labels.csv
go run . -t "Col: MultiTimeline" -c column test/multiTimeline.csv
go run . -t "Col: Single Column Header" -c column test/single_col_header.csv
go run . -t "Col: Single Column" -c column test/single_col.csv
go run . -t "Col: Two Columns" -c column test/two_col.csv
go run . -t "Col: Two Columns Header" -c column test/two_col_header.csv