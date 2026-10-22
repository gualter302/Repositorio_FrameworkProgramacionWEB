$ErrorActionPreference = 'Stop'
$base = 'http://localhost:3008/api'
$session = Invoke-RestMethod -Method Post -Uri "$base/login" -ContentType 'application/json' -Body '{"email":"cliente@upse.edu.ec","password":"123456"}'
$headers = @{ Authorization = "Bearer $($session.token)"; 'Idempotency-Key' = [guid]::NewGuid().ToString() }
$body = '{"cliente":{"nombre":"Prueba PWA","email":"cliente@upse.edu.ec","direccion":"Centro","ciudad":"La Libertad","telefono":"0999999999"},"metodoPago":"Efectivo","items":[{"productoId":101,"cantidad":1}]}'
$before = Invoke-RestMethod "$base/productos/101"
$first = Invoke-RestMethod -Method Post -Uri "$base/checkout" -Headers $headers -ContentType 'application/json' -Body $body
$second = Invoke-RestMethod -Method Post -Uri "$base/checkout" -Headers $headers -ContentType 'application/json' -Body $body
$after = Invoke-RestMethod "$base/productos/101"
if ($first.id -ne $second.id -or $after.stock -ne ($before.stock - 1)) { throw 'El reintento duplicó la compra o el descuento de stock' }
try {
    Invoke-RestMethod -Method Post -Uri "$base/checkout" -Headers $headers -ContentType 'application/json' -Body ($body.Replace('"cantidad":1','"cantidad":2')) | Out-Null
    throw 'Se aceptó una clave con otro carrito'
} catch {
    if ([int]$_.Exception.Response.StatusCode -ne 409) { throw }
}
Write-Output "OK: pedido $($first.numero) recuperado sin duplicar stock; cambio de carrito rechazado con 409."
$headers['Idempotency-Key'] = [guid]::NewGuid().ToString()
try {
    Invoke-RestMethod -Method Post -Uri "$base/checkout" -Headers $headers -ContentType 'application/json' -Body ($body.Replace('"cantidad":1','"cantidad":99999')) | Out-Null
    throw 'Se aceptó una compra sin stock'
} catch {
    if ([int]$_.Exception.Response.StatusCode -ne 409) { throw }
}
$stockFinal = Invoke-RestMethod "$base/productos/101"
if ($stockFinal.stock -ne $after.stock) { throw 'La compra fallida modificó stock' }
Write-Output 'OK: stock insuficiente responde 409 y conserva existencias.'
