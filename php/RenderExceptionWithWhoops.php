<?php

namespace App\Traits;

use Whoops\Handler\JsonResponseHandler;
use Whoops\Handler\PrettyPageHandler;
use Whoops\Util\Misc;

trait RenderExceptionWithWhoops
{
    private function renderExceptionWithWhoops(\Exception $e)
    {
        $handler = Misc::isAjaxRequest() ? new JsonResponseHandler : new PrettyPageHandler;
        if ($editor = config('whoops.editor')) {
            // https://github.com/filp/whoops/blob/master/docs/Open%20Files%20In%20An%20Editor.md
            $handler->setEditor('sublime');
        }
        
        $whoops = new \Whoops\Run;
        $whoops->pushHandler($handler);
        return new \Illuminate\Http\Response(
            $whoops->handleException($e),
            $e->getStatusCode(),
            $e->getHeaders()
        );
    }
}
