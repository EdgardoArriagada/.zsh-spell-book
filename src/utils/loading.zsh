loading() {
  local seconds=${1?Error: Amount of seconds expected.}
  local timeFragment=$(($seconds / 25.0))

  printf '[>                    ]   (0%%)\r'
  sleep "$timeFragment"
  printf '[=>                   ]   (5%%)\r'
  sleep "$timeFragment"
  printf '[==>                  ]   (10%%)\r'
  sleep "$timeFragment"
  printf '[===>                 ]   (15%%)\r'
  sleep "$timeFragment"
  printf '[====>                ]   (20%%)\r'
  sleep "$timeFragment"
  printf '[=====>               ]   (25%%)\r'
  sleep "$timeFragment"
  printf '[======>              ]   (30%%)\r'
  sleep "$timeFragment"
  printf '[=======>             ]   (35%%)\r'
  sleep "$timeFragment"
  printf '[========>            ]   (40%%)\r'
  sleep "$timeFragment"
  printf '[=========>           ]   (45%%)\r'
  sleep "$timeFragment"
  printf '[==========>          ]   (50%%)\r'
  sleep "$timeFragment"
  printf '[===========>         ]   (55%%)\r'
  sleep "$timeFragment"
  printf '[============>        ]   (60%%)\r'
  sleep "$timeFragment"
  printf '[=============>       ]   (65%%)\r'
  sleep "$timeFragment"
  printf '[==============>      ]   (70%%)\r'
  sleep "$timeFragment"
  printf '[===============>     ]   (75%%)\r'
  sleep "$timeFragment"
  printf '[================>    ]   (80%%)\r'
  sleep "$timeFragment"
  printf '[=================>   ]   (85%%)\r'
  sleep "$timeFragment"
  printf '[==================>  ]   (90%%)\r'
  sleep "$timeFragment"
  printf '[===================> ]   (95%%)\r'
  sleep "$timeFragment"
  printf '[====================>]   (100%%)\r\n'
}
