angular.module('cells', ['ngResource', 'angularTreeview'])
  .factory('Cells', function($resource) {
    var Cells = $resource('/v1/cells',{},{
                  query: {
                    method: 'GET',
                    isArray: true,
                    headers: {
                      'X-API-Token': 'bla'
                    },
                  }
                });


    return Cells;
  })
  .factory('CellFull', function($resource) {
    var CellFull = $resource('/v1/cell/:id/full',{id: '99'},{
                  get: {
                    method: 'GET',
                    isArray: false,
                    headers: {
                      'X-API-Token': 'bla'
                    },
                  }
                });

    return CellFull;
  })
  .controller('CellsCtrl', function($scope, Cells) {
    $scope.cells = Cells.query();
    //console.log($scope.cells)
  })
  .controller('CellFullCtrl', function($scope, CellFull) {
    $scope._full_cell = CellFull.get(function(res) {
      $scope.full_cell = []
      var cell_root = {}
      cell_root.label = 'cell'
      cell_root.children = []
      angular.forEach(JSON.parse(angular.toJson(res)), function(value, key) {
        if ( angular.isObject(value) ) {
          cell_root.children.push(buildNode(key, value))

        } else {
          cell_root[key] = value
        }
      })

      $scope.full_cell.push(cell_root)
      console.log(cell_root)
    });

    $scope.$watch('cellInfos.currentNode', function(newObj, oldObj) {
      if ( $scope.cellInfos && angular.isObject($scope.cellInfos.currentNode) ) {
        $scope.formData = $scope.cellInfos.currentNode
        console.log( $scope.cellInfos.currentNode )
      }
    }, false);
  });

function buildNode(key, value) {
  var node = {}
  node.children = []
  node.label = key

  angular.forEach(value, function(v, k) {
    if ( angular.isArray(value) ) {
      node.children.push(buildNode(key, v))

    } else if ( angular.isObject(v) ) {
      node.children.push(buildNode(k, v))
    } else {
      node[k] = v
    }
  })
  if ( !node.name ) {
    node.name = key
  }

  return node
}
