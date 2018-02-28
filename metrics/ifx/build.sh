set -x
rm -f libifx.so

CXXFLAGS="-I. -I/usr/include/azurepal/IfxMetrics -I/usr/include/azurepal -I/usr/include/azurepal/etw -I/usr/include -I/usr/include/azurepal/unix/winapi -I/usr/include/azurepal/unix/winapi/crt -I/usr/include/c++/4.8 -I/usr/include/x86_64-linux-gnu/c++/4.8"
LDFLAGS="-L/usr/lib/x86_64-linux-gnu -lrtcpal -lIfxMetrics"
g++ -fPIC -shared -o libifx.so ifx.cpp -lstdc++ -std=c++11 ${CXXFLAGS} ${LDFLAGS}
cp libifx.so /usr/lib/x86_64-linux-gnu/libifx.so