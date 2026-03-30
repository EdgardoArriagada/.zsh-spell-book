..() { builtin cd .. && lsd; }
...() { builtin cd ../.. && lsd; }
....() { builtin cd ../../.. && lsd; }
.....() { builtin cd ../../../.. && lsd; }
......() { builtin cd ../../../../.. && lsd; }
.......() { builtin cd ../../../../../.. && lsd; }

hisIgnore '..'
