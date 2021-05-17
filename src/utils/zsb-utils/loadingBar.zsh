${zsb}.loadingBar() {
  local -r seconds=${1:?'Amount of seconds expecte.'}
  local timeFragment=$(($seconds / 20.0))

  echo -ne '[>                    ]   (0%)\r'
  sleep "$timeFragment"
  sleep "$timeFragment"
  echo -ne '[=>                   ]   (5%)\r'
  sleep "$timeFragment"
  echo -ne '[==>                  ]   (10%)\r'
  sleep "$timeFragment"
  echo -ne '[===>                 ]   (15%)\r'
  sleep "$timeFragment"
  echo -ne '[====>                ]   (20%)\r'
  sleep "$timeFragment"
  echo -ne '[=====>               ]   (25%)\r'
  sleep "$timeFragment"
  echo -ne '[======>              ]   (30%)\r'
  sleep "$timeFragment"
  echo -ne '[=======>             ]   (35%)\r'
  sleep "$timeFragment"
  echo -ne '[========>            ]   (40%)\r'
  sleep "$timeFragment"
  echo -ne '[=========>           ]   (45%)\r'
  sleep "$timeFragment"
  echo -ne '[==========>          ]   (50%)\r'
  sleep "$timeFragment"
  echo -ne '[===========>         ]   (55%)\r'
  sleep "$timeFragment"
  echo -ne '[============>        ]   (60%)\r'
  sleep "$timeFragment"
  echo -ne '[=============>       ]   (65%)\r'
  sleep "$timeFragment"
  echo -ne '[==============>      ]   (70%)\r'
  sleep "$timeFragment"
  echo -ne '[===============>     ]   (75%)\r'
  sleep "$timeFragment"
  echo -ne '[================>    ]   (80%)\r'
  sleep "$timeFragment"
  echo -ne '[=================>   ]   (85%)\r'
  sleep "$timeFragment"
  echo -ne '[==================>  ]   (90%)\r'
  sleep "$timeFragment"
  echo -ne '[===================> ]   (95%)\r'
  sleep "$timeFragment"
  echo -ne '[====================>]   (100%)\r'
  echo -ne '\n'
}
