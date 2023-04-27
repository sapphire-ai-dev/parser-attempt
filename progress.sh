total=$(find . -name '*.go' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
test=$(find . -name '*test.go' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
data=$(find . -name '*.txt' | sed 's/.*/"&"/' | xargs  wc -l | tail -1 | awk '{ print $1 }')
code=$((total-test))
total=$((code + test + data))
# shellcheck disable=SC2004
echo "impl: $code, test: $test, data: $data, total: $total"
